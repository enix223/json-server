package serializer

type Serializer interface {
	Serialize(payload interface{}) ([]byte, error)
}
