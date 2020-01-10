package v1alpha1

import (
	"bufio"
	"bytes"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	istiov1alpha1 "istio.io/api/rbac/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleBinding is a Istio ServiceRoleBinding resource
type ServiceRoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceRoleBindingSpec `json:"spec"`
}

func (vs *ServiceRoleBinding) GetSpecMessage() proto.Message {
	return &vs.Spec.ServiceRoleBinding
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleBindingList is a list of ServiceRoleBinding resources
type ServiceRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ServiceRoleBinding `json:"items"`
}

// ServiceRoleBindingSpec is a wrapper around Istio ServiceRole
type ServiceRoleBindingSpec struct {
	istiov1alpha1.ServiceRoleBinding
}

func (p *ServiceRoleBindingSpec) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(writer, &p.ServiceRoleBinding)
	if err != nil {
		log.Printf("Could not marshal ServiceRoleBindingSpec. Error: %v", err)
		return nil, err
	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func (p *ServiceRoleBindingSpec) UnmarshalJSON(b []byte) error {
	reader := bytes.NewReader(b)
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(reader, &p.ServiceRoleBinding)
	if err != nil {
		log.Printf("Could not unmarshal ServiceRoleBindingSpec. Error: %v", err)
		return err
	}
	return nil
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceRoleBindingSpec) DeepCopyInto(out *ServiceRoleBindingSpec) {
	*out = *in
}
