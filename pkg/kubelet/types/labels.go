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

package types

const (
	KubernetesPodNameLabel       = "io.kubernetes.pod.name"
	KubernetesPodNamespaceLabel  = "io.kubernetes.pod.namespace"
	KubernetesPodTenantLabel     = "io.kubernetes.pod.tenant"
	KubernetesPodUIDLabel        = "io.kubernetes.pod.uid"
	KubernetesContainerNameLabel = "io.kubernetes.container.name"
)

func GetContainerName(labels map[string]string) string {
	return labels[KubernetesContainerNameLabel]
}

func GetPodName(labels map[string]string) string {
	return labels[KubernetesPodNameLabel]
}

func GetPodUID(labels map[string]string) string {
	return labels[KubernetesPodUIDLabel]
}

func GetPodNamespace(labels map[string]string) string {
	return labels[KubernetesPodNamespaceLabel]
}

func GetPodTenant(labels map[string]string) string {
	return labels[KubernetesPodTenantLabel]
}