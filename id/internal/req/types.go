package req

type Body interface {
	Validate() error
}

type Request[T any] struct {
	Meta *Meta `json:"meta,omitempty"`
	Data T     `json:"data"`
}

type Meta struct{}
