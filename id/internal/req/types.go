package req

type Body interface {
	Validate() error
}

type Request[T any] struct {
	Meta *RequestMeta `json:"meta,omitempty"`
	Data T            `json:"data,omitempty"`
}

type RequestMeta struct{}
