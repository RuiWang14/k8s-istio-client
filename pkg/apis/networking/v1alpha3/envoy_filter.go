package v1alpha3

import (
	"bufio"
	"bytes"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EnvoyFilter is an Istio EnvoyFilter resource
type EnvoyFilter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EnvoyFilterSpec `json:"spec"`
}

func (vs *EnvoyFilter) GetSpecMessage() proto.Message {
	return &vs.Spec.EnvoyFilter
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EnvoyFilterList is a list of EnvoyFilter resources
type EnvoyFilterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []EnvoyFilter `json:"items"`
}

// EnvoyFilterSpec is a wrapper around Istio EnvoyFilter
type EnvoyFilterSpec struct {
	istiov1alpha3.EnvoyFilter
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvoyFilterSpec) DeepCopyInto(out *EnvoyFilterSpec) {
	*out = *in
}

func (vs *EnvoyFilterSpec) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(writer, &vs.EnvoyFilter)
	if err != nil {
		return nil, err
	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func (vs *EnvoyFilterSpec) UnmarshalJSON(b []byte) error {
	reader := bytes.NewReader(b)
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(reader, &vs.EnvoyFilter)
	if err != nil {
		return err
	}
	return nil
}
