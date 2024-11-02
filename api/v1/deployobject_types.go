/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeployObjectSpec defines the desired state of DeployObject
// Kubernetes 通过将期望状态(Spec) 与实际集群状态（其他对象的 Status）和外部状态进行协调，然后记录它观察到的内容 （ Status ） 来发挥作用。
// 因此，每个功能对象都包括规范和状态
// spec 保存预期状态，因此控制器的任何“输入”都在这里
type DeployObjectSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DeployObject. Edit deployobject_types.go to remove/update
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:ExclusiveMaximum=false
	Size int32 `json:"size,omitempty"` //所有序列化字段都必须为 camelCase ，因此我们使用 JSON 结构标签来指定这一点
}

// DeployObjectStatus defines the observed state of DeployObject
type DeployObjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DeployObject is the Schema for the deployobjects API
// CronJob 是我们的根类型，描述了 CronJob 种。像所有 Kubernetes 对象一样，它包含 TypeMeta （描述 API 版本和种类）
// 还包含 ObjectMeta 个，其中包含名称、命名空间和标签等内容。CronJobList 只是多个 CronJob 的容器.
type DeployObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeployObjectSpec   `json:"spec,omitempty"`
	Status DeployObjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeployObjectList contains a list of DeployObject
type DeployObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeployObject `json:"items"`
}

// 将 Go 类型添加到 API 组。这允许我们将此 API 组中的类型添加到任何 scheme 中
func init() {
	SchemeBuilder.Register(&DeployObject{}, &DeployObjectList{})
}
