package serializer

import (
	"fmt"
)

type textSerializer struct {
}

// Serialize implements Serializer
func (j *textSerializer) Serialize(payload interface{}) ([]byte, error) {
	v, ok := payload.(string)
	if !ok {
		return nil, fmt.Errorf("failed to serialize: %v", payload)
	}
	return []byte(v), nil
}

func NewTextSerializer() Serializer {
	return &textSerializer{}
}
