package serializer

import "encoding/json"

type jsonSerializer struct {
}

// Serialize implements Serializer
func (j *jsonSerializer) Serialize(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

func NewJsonSerializer() Serializer {
	return &jsonSerializer{}
}
