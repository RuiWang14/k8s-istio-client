package v1alpha1

import (
	"bufio"
	"bytes"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	istiov1alpha1 "istio.io/api/authentication/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MeshPolicy is a Istio MeshPolicy resource
type MeshPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MeshPolicySpec `json:"spec"`
}

// GetSpecMessage gets the Policy in the MeshPolicy Spec
func (p *MeshPolicy) GetSpecMessage() proto.Message {
	return &p.Spec.Policy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MeshPolicyList is a list of MeshPolicy resources
type MeshPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MeshPolicy `json:"items"`
}

// MeshPolicySpec in a wrapper around Istio Policy
type MeshPolicySpec struct {
	istiov1alpha1.Policy
}

func (p *MeshPolicySpec) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(writer, &p.Policy)
	if err != nil {
		log.Printf("Could not marshal PolicySpec. Error: %v", err)
		return nil, err
	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func (p *MeshPolicySpec) UnmarshalJSON(b []byte) error {
	reader := bytes.NewReader(b)
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(reader, &p.Policy)
	if err != nil {
		log.Printf("Could not unmarshal PolicySpec. Error: %v", err)
		return err
	}
	return nil
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeshPolicySpec) DeepCopyInto(out *MeshPolicySpec) {
	*out = *in
}
