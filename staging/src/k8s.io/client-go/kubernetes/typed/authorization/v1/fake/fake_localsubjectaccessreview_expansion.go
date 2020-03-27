/*
Copyright 2017 The Kubernetes Authors.
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

package fake

import (
	authorizationapi "k8s.io/api/authorization/v1"
	core "k8s.io/client-go/testing"
)

func (c *FakeLocalSubjectAccessReviews) Create(sar *authorizationapi.LocalSubjectAccessReview) (result *authorizationapi.LocalSubjectAccessReview, err error) {
	obj, err := c.Fake.Invokes(core.NewCreateActionWithMultiTenancy(authorizationapi.SchemeGroupVersion.WithResource("localsubjectaccessreviews"), c.ns, sar, c.te), &authorizationapi.SubjectAccessReview{})
	return obj.(*authorizationapi.LocalSubjectAccessReview), err
}