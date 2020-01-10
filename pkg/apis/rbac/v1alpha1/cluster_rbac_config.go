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

// ClusterRbacConfig is a Istio ClusterRbacConfig resource
type ClusterRbacConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterRbacConfigSpec `json:"spec"`
}

func (vs *ClusterRbacConfig) GetSpecMessage() proto.Message {
	return &vs.Spec.RbacConfig
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRbacConfigList is a list of ClusterRbacConfig resources
type ClusterRbacConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterRbacConfig `json:"items"`
}

// ClusterRbacConfigSpec is a wrapper around Istio ClusterRbacConfig
type ClusterRbacConfigSpec struct {
	istiov1alpha1.RbacConfig
}

func (p *ClusterRbacConfigSpec) MarshalJSON() ([]byte, error) {
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

func (p *ClusterRbacConfigSpec) UnmarshalJSON(b []byte) error {
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
func (in *ClusterRbacConfigSpec) DeepCopyInto(out *ClusterRbacConfigSpec) {
	*out = *in
}
