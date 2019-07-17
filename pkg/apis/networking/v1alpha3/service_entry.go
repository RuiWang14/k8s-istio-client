package v1alpha3

import (
	"bufio"
	"bytes"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceEntry is a Istio ServiceEntry resource
type ServiceEntry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceEntrySpec `json:"spec"`
}

func (vs *ServiceEntry) GetSpecMessage() proto.Message {
	return &vs.Spec.ServiceEntry
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceEntryList is a list of ServiceEntry resources
type ServiceEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ServiceEntry `json:"items"`
}

// ServiceEntrySpec is a wrapper around Istio ServiceEntry
type ServiceEntrySpec struct {
	istiov1alpha3.ServiceEntry
}

func (p *ServiceEntrySpec) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(writer, &p.ServiceEntry)
	if err != nil {
		log.Printf("Could not marshal ServiceEntrySpec. Error: %v", err)
		return nil, err
	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func (p *ServiceEntrySpec) UnmarshalJSON(b []byte) error {
	reader := bytes.NewReader(b)
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(reader, &p.ServiceEntry)
	if err != nil {
		log.Printf("Could not unmarshal ServiceEntrySpec. Error: %v", err)
		return err
	}
	return nil
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceEntrySpec) DeepCopyInto(out *ServiceEntrySpec) {
	*out = *in
}
