package res

type Response[T any] struct {
	Meta  *ResponseMeta  `json:"meta"`
	Data  T              `json:"data"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ResponseMeta struct{}

type ErrorResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
