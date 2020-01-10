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

// RbacConfig is a Istio RbacConfig resource
type RbacConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RbacConfigSpec `json:"spec"`
}

func (vs *RbacConfig) GetSpecMessage() proto.Message {
	return &vs.Spec.RbacConfig
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RbacConfigList is a list of RbacConfig resources
type RbacConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RbacConfig `json:"items"`
}

// RbacConfigSpec is a wrapper around Istio RbacConfig
type RbacConfigSpec struct {
	istiov1alpha1.RbacConfig
}

func (p *RbacConfigSpec) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(writer, &p.RbacConfig)
	if err != nil {
		log.Printf("Could not marshal RbacConfigSpec. Error: %v", err)
		return nil, err
	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func (p *RbacConfigSpec) UnmarshalJSON(b []byte) error {
	reader := bytes.NewReader(b)
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(reader, &p.RbacConfig)
	if err != nil {
		log.Printf("Could not unmarshal RbacConfigSpec. Error: %v", err)
		return err
	}
	return nil
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RbacConfigSpec) DeepCopyInto(out *RbacConfigSpec) {
	*out = *in
}
