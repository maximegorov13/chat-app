package req

type Body interface {
	Validate() error
}

type Request[T any] struct {
	Meta *RequestMeta `json:"meta,omitempty"`
	Data T            `json:"data"`
}

type RequestMeta struct{}
