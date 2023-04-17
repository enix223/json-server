package server

type PayloadType string

const (
	PayloadTypeJson PayloadType = "json"
	PayloadTypeText PayloadType = "text"
)

type EndpointMapping struct {
	Path            string             `json:"path"`
	Method          *string            `json:"method"`
	Payload         *interface{}       `json:"payload"`
	PayloadType     PayloadType        `json:"payload_type"`
	ResponseHeaders *map[string]string `json:"resposne_headers"`
	StatusCode      *int               `json:"status_code"`
}
