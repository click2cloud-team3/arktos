/*
Copyright 2016 The Kubernetes Authors.
Copyright 2020 Authors of Arktos - file modified.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package etcd3

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"

	apitesting "k8s.io/apimachinery/pkg/api/apitesting"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/apis/example"
	examplev1 "k8s.io/apiserver/pkg/apis/example/v1"
	"k8s.io/apiserver/pkg/storage"
)

func TestWatch(t *testing.T) {
	testWatch(t, false)
}

func TestWatchList(t *testing.T) {
	testWatch(t, true)
}

// It tests that
// - first occurrence of objects should notify Add event
// - update should trigger Modified event
// - update that gets filtered should trigger Deleted event
func testWatch(t *testing.T, recursive bool) {
	ctx, store, cluster := testSetup(t)
	defer cluster.Terminate(t)
	podFoo := &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}
	podBar := &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "bar"}}

	tests := []struct {
		key        string
		pred       storage.SelectionPredicate
		watchTests []*testWatchStruct
	}{{ // create a key
		key:        "/somekey-1",
		watchTests: []*testWatchStruct{{podFoo, true, watch.Added}},
		pred:       storage.Everything,
	}, { // create a key but obj gets filtered. Then update it with unfiltered obj
		key:        "/somekey-3",
		watchTests: []*testWatchStruct{{podFoo, false, ""}, {podBar, true, watch.Added}},
		pred: storage.SelectionPredicate{
			Label: labels.Everything(),
			Field: fields.ParseSelectorOrDie("metadata.name=bar"),
			GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
				pod := obj.(*example.Pod)
				return nil, fields.Set{"metadata.name": pod.Name}, nil
			},
		},
	}, { // update
		key:        "/somekey-4",
		watchTests: []*testWatchStruct{{podFoo, true, watch.Added}, {podBar, true, watch.Modified}},
		pred:       storage.Everything,
	}, { // delete because of being filtered
		key:        "/somekey-5",
		watchTests: []*testWatchStruct{{podFoo, true, watch.Added}, {podBar, true, watch.Deleted}},
		pred: storage.SelectionPredicate{
			Label: labels.Everything(),
			Field: fields.ParseSelectorOrDie("metadata.name!=bar"),
			GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
				pod := obj.(*example.Pod)
				return nil, fields.Set{"metadata.name": pod.Name}, nil
			},
		},
	}}
	for i, tt := range tests {
		aw := store.watch(ctx, tt.key, "0", tt.pred, recursive)
		if aw.GetFirstError() != nil {
			t.Fatalf("Watch failed: %v", aw.GetFirstError())
		}
		var prevObj *example.Pod
		for _, watchTest := range tt.watchTests {
			out := &example.Pod{}
			key := tt.key
			if recursive {
				key = key + "/item"
			}
			err := store.GuaranteedUpdate(ctx, key, out, true, nil, storage.SimpleUpdate(
				func(runtime.Object) (runtime.Object, error) {
					return watchTest.obj, nil
				}))
			if err != nil {
				t.Fatalf("GuaranteedUpdate failed: %v", err)
			}
			if watchTest.expectEvent {
				expectObj := out
				if watchTest.watchType == watch.Deleted {
					expectObj = prevObj
					expectObj.ResourceVersion = out.ResourceVersion
				}
				testCheckResult(t, i, watchTest.watchType, aw, expectObj)
			}
			prevObj = out
		}
		aw.Stop()
		testCheckStop(t, i, aw)
	}
}

func TestDeleteTriggerWatch(t *testing.T) {
	ctx, store, cluster := testSetup(t)
	defer cluster.Terminate(t)
	key, storedObj := testPropogateStore(ctx, t, store, &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})
	aw := store.Watch(ctx, key, storedObj.ResourceVersion, storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatalf("Watch failed: %v", aw.GetFirstError())
	}
	if err := store.Delete(ctx, key, &example.Pod{}, nil, storage.ValidateAllObjectFunc); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	testCheckEventType(t, watch.Deleted, aw)
}

// TestWatchFromZero tests that
// - watch from 0 should sync up and grab the object added before
// - watch from 0 is able to return events for objects whose previous version has been compacted
func TestWatchFromZero(t *testing.T) {
	ctx, store, cluster := testSetup(t)
	defer cluster.Terminate(t)
	key, storedObj := testPropogateStore(ctx, t, store, &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns"}})

	aw := store.Watch(ctx, key, "0", storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatalf("Watch failed: %v", aw.GetFirstError())
	}
	testCheckResult(t, 0, watch.Added, aw, storedObj)
	aw.Stop()

	// Update
	out := &example.Pod{}
	err := store.GuaranteedUpdate(ctx, key, out, true, nil, storage.SimpleUpdate(
		func(runtime.Object) (runtime.Object, error) {
			return &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns", Annotations: map[string]string{"a": "1"}}}, nil
		}))
	if err != nil {
		t.Fatalf("GuaranteedUpdate failed: %v", err)
	}

	// Make sure when we watch from 0 we receive an ADDED event
	aw = store.Watch(ctx, key, "0", storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatalf("Watch failed: %v", aw.GetFirstError())
	}
	testCheckResult(t, 1, watch.Added, aw, out)
	aw.Stop()

	// Update again
	out = &example.Pod{}
	err = store.GuaranteedUpdate(ctx, key, out, true, nil, storage.SimpleUpdate(
		func(runtime.Object) (runtime.Object, error) {
			return &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns"}}, nil
		}))
	if err != nil {
		t.Fatalf("GuaranteedUpdate failed: %v", err)
	}

	// Compact previous versions
	revToCompact, err := store.versioner.ParseResourceVersion(out.ResourceVersion)
	if err != nil {
		t.Fatalf("Error converting %q to an int: %v", storedObj.ResourceVersion, err)
	}
	_, err = cluster.RandClient().Compact(ctx, int64(revToCompact), clientv3.WithCompactPhysical())
	if err != nil {
		t.Fatalf("Error compacting: %v", err)
	}

	// Make sure we can still watch from 0 and receive an ADDED event
	aw = store.Watch(ctx, key, "0", storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatalf("Watch failed: %v", aw.GetFirstError())
	}
	testCheckResult(t, 2, watch.Added, aw, out)
}

// TestWatchFromNoneZero tests that
// - watch from non-0 should just watch changes after given version
func TestWatchFromNoneZero(t *testing.T) {
	ctx, store, cluster := testSetup(t)
	defer cluster.Terminate(t)
	key, storedObj := testPropogateStore(ctx, t, store, &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})

	aw := store.Watch(ctx, key, storedObj.ResourceVersion, storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatalf("Watch failed: %v", aw.GetFirstError())
	}
	out := &example.Pod{}
	store.GuaranteedUpdate(ctx, key, out, true, nil, storage.SimpleUpdate(
		func(runtime.Object) (runtime.Object, error) {
			return &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "bar"}}, aw.GetFirstError()
		}))
	testCheckResult(t, 0, watch.Modified, aw, out)
}

func TestWatchError(t *testing.T) {
	codec := &testCodec{apitesting.TestCodec(codecs, examplev1.SchemeGroupVersion)}
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer cluster.Terminate(t)
	invalidStore := newStore(cluster.RandClient(), true, codec, "", prefixTransformer{prefix: []byte("test!")})
	ctx := context.Background()
	aw := invalidStore.Watch(ctx, "/abc", "0", storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatalf("Watch failed: %v", aw.GetFirstError())
	}
	validStore := newStore(cluster.RandClient(), true, codec, "", prefixTransformer{prefix: []byte("test!")})
	validStore.GuaranteedUpdate(ctx, "/abc", &example.Pod{}, true, nil, storage.SimpleUpdate(
		func(runtime.Object) (runtime.Object, error) {
			return &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}, nil
		}))
	testCheckEventType(t, watch.Error, aw)
}

func TestWatchContextCancel(t *testing.T) {
	ctx, store, cluster := testSetup(t)
	defer cluster.Terminate(t)
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()
	// When we watch with a canceled context, we should detect that it's context canceled.
	// We won't take it as error and also close the watcher.
	aw := store.watcher.Watch(canceledCtx, "/abc", 0, false, storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatal(aw.GetFirstError())
	}

	select {
	case _, ok := <-aw.ResultChan():
		if ok {
			t.Error("ResultChan() should be closed")
		}
	case <-time.After(wait.ForeverTestTimeout):
		t.Errorf("timeout after %v", wait.ForeverTestTimeout)
	}
}

func TestWatchErrResultNotBlockAfterCancel(t *testing.T) {
	origCtx, store, cluster := testSetup(t)
	defer cluster.Terminate(t)
	ctx, cancel := context.WithCancel(origCtx)
	w := store.watcher.createWatchChan(ctx, "/abc", 0, false, storage.Everything)
	// make resutlChan and errChan blocking to ensure ordering.
	w.resultChan = make(chan watch.Event)
	w.errChan = make(chan error)
	// The event flow goes like:
	// - first we send an error, it should block on resultChan.
	// - Then we cancel ctx. The blocking on resultChan should be freed up
	//   and run() goroutine should return.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w.run()
		wg.Done()
	}()
	w.errChan <- fmt.Errorf("some error")
	cancel()
	wg.Wait()
}

func TestWatchDeleteEventObjectHaveLatestRV(t *testing.T) {
	ctx, store, cluster := testSetup(t)
	defer cluster.Terminate(t)
	key, storedObj := testPropogateStore(ctx, t, store, &example.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})

	aw := store.Watch(ctx, key, storedObj.ResourceVersion, storage.Everything)
	if aw.GetFirstError() != nil {
		t.Fatalf("Watch failed: %v", aw.GetFirstError())
	}
	etcdW := cluster.RandClient().Watch(ctx, "/", clientv3.WithPrefix())

	if err := store.Delete(ctx, key, &example.Pod{}, &storage.Preconditions{}, storage.ValidateAllObjectFunc); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	e := <-aw.ResultChan()
	watchedDeleteObj := e.Object.(*example.Pod)
	var wres clientv3.WatchResponse
	wres = <-etcdW

	watchedDeleteRev, err := store.versioner.ParseResourceVersion(watchedDeleteObj.ResourceVersion)
	if err != nil {
		t.Fatalf("ParseWatchResourceVersion failed: %v", err)
	}
	if int64(watchedDeleteRev) != wres.Events[0].Kv.ModRevision {
		t.Errorf("Object from delete event have version: %v, should be the same as etcd delete's mod rev: %d",
			watchedDeleteRev, wres.Events[0].Kv.ModRevision)
	}
}

type testWatchStruct struct {
	obj         *example.Pod
	expectEvent bool
	watchType   watch.EventType
}

type testCodec struct {
	runtime.Codec
}

func (c *testCodec) Decode(data []byte, defaults *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	return nil, nil, errTestingDecode
}

func testCheckEventType(t *testing.T, expectEventType watch.EventType, w watch.Interface) {
	select {
	case res := <-w.ResultChan():
		if res.Type != expectEventType {
			t.Errorf("event type want=%v, get=%v", expectEventType, res.Type)
		}
	case <-time.After(wait.ForeverTestTimeout):
		t.Errorf("time out after waiting %v on ResultChan", wait.ForeverTestTimeout)
	}
}

func testCheckResult(t *testing.T, i int, expectEventType watch.EventType, w watch.Interface, expectObj *example.Pod) {
	select {
	case res := <-w.ResultChan():
		if res.Type != expectEventType {
			t.Errorf("#%d: event type want=%v, get=%v", i, expectEventType, res.Type)
			return
		}
		if !reflect.DeepEqual(expectObj, res.Object) {
			t.Errorf("#%d: obj want=\n%#v\nget=\n%#v", i, expectObj, res.Object)
		}
	case <-time.After(wait.ForeverTestTimeout):
		t.Errorf("#%d: time out after waiting %v on ResultChan", i, wait.ForeverTestTimeout)
	}
}

func testCheckStop(t *testing.T, i int, w watch.Interface) {
	select {
	case e, ok := <-w.ResultChan():
		if ok {
			var obj string
			switch e.Object.(type) {
			case *example.Pod:
				obj = e.Object.(*example.Pod).Name
			case *metav1.Status:
				obj = e.Object.(*metav1.Status).Message
			}
			t.Errorf("#%d: ResultChan should have been closed. Event: %s. Object: %s", i, e.Type, obj)
		}
	case <-time.After(wait.ForeverTestTimeout):
		t.Errorf("#%d: time out after waiting 1s on ResultChan", i)
	}
}

func TestGetKeyAndOptFromPartitionConfig(t *testing.T) {
	for _, tc := range []struct {
		name               string
		key                string
		partitionedConfig  map[string]storage.Interval
		expectedKey        string
		expectedOpt        clientv3.OpOption
	}{
		{
			name: "opt with interval beginning and end via key",
			key:  "registry/pods/",
			partitionedConfig: map[string]storage.Interval{
				"registry/pods/": storage.Interval{
					Begin: "tenant1",
					End:   "tenant2",
				},
				"registry/pod/": storage.Interval{
					Begin: "tenant3",
					End:   "tenant4",
				},
			},
			expectedKey: "registry/pods/tenant1",
			expectedOpt: clientv3.WithRange("registry/pods/tenant2"),
		},
		{
			name: "opt with interval beginning only via key",
			key:  "registry/deployment/",
			partitionedConfig: map[string]storage.Interval{
				"registry/deployment/": storage.Interval{
					Begin: "tenant1",
				},
				"registry/pod/": storage.Interval{
					Begin: "tenant3",
					End:   "tenant4",
				},
			},
			expectedKey: "registry/deployment/tenant1",
			expectedOpt: clientv3.WithFromKey(),
		},
		{
			name: "opt with interval end only via key",
			key:  "registry/pods/",
			partitionedConfig: map[string]storage.Interval{
				"registry/pods/": storage.Interval{
					End: "tenant",
				},
				"registry/pod/": storage.Interval{
					Begin: "tenant3",
					End:   "tenant4",
				},
			},
			expectedKey: "registry/pods/",
			expectedOpt: clientv3.WithRange("tenant"),
		},
		{
			name: "opt without beginning and end",
			key:  "registry/pods/",
			partitionedConfig: map[string]storage.Interval{
				"registry/po/": storage.Interval{
					End: "tenant",
				},
				"registry/pod/": storage.Interval{
					Begin: "tenant3",
					End:   "tenant4",
				},
			},
			expectedKey: "registry/pods/",
			expectedOpt: nil,
		},
		{
			name: "opt with empty entry via key",
			key:  "registry/replicaset/",
			partitionedConfig: map[string]storage.Interval{
				"registry/replicaset/": storage.Interval{
				},
				"registry/pod/": storage.Interval{
					Begin: "tenant3",
					End:   "tenant4",
				},
			},
			expectedKey: "registry/replicaset/",
			expectedOpt: nil,
		},
		{
			name: "opt with empty config",
			key:  "registry/demonset/",
			partitionedConfig: map[string]storage.Interval{
			},
			expectedKey: "registry/demonset/",
			expectedOpt: nil,
		},

	} {
		t.Run(tc.name, func(t *testing.T) {
			updatedKey, opt := GetKeyAndOptFromPartitionConfig( tc.key, tc.partitionedConfig)
			if tc.expectedKey != updatedKey {
				t.Fatalf("The key %s is not expected as %s", updatedKey, tc.expectedKey)
			}
			if reflect.ValueOf(tc.expectedOpt).Pointer() != reflect.ValueOf(opt).Pointer() {
				t.Fatalf("The opt %s is not expected as %s ", reflect.ValueOf(opt).String(), reflect.ValueOf(tc.expectedOpt).String())
			}
		})
	}
}
