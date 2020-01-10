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

// ServiceRole is a Istio ServiceRole resource
type ServiceRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceRoleSpec `json:"spec"`
}

func (vs *ServiceRole) GetSpecMessage() proto.Message {
	return &vs.Spec.ServiceRole
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleList is a list of ServiceRole resources
type ServiceRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ServiceRole `json:"items"`
}

// ServiceRoleSpec is a wrapper around Istio ServiceRole
type ServiceRoleSpec struct {
	istiov1alpha1.ServiceRole
}

func (p *ServiceRoleSpec) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(writer, &p.ServiceRole)
	if err != nil {
		log.Printf("Could not marshal ServiceRoleSpec. Error: %v", err)
		return nil, err
	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func (p *ServiceRoleSpec) UnmarshalJSON(b []byte) error {
	reader := bytes.NewReader(b)
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(reader, &p.ServiceRole)
	if err != nil {
		log.Printf("Could not unmarshal ServiceRoleSpec. Error: %v", err)
		return err
	}
	return nil
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceRoleSpec) DeepCopyInto(out *ServiceRoleSpec) {
	*out = *in
}
