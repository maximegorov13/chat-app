package res

type Response struct {
	Meta  Meta          `json:"meta"`
	Data  any           `json:"data,omitempty"`
	Error *ErrorDetails `json:"error,omitempty"`
}

type Meta struct{}

type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
