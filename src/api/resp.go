package api

type ApiResponse struct {
	TraceID string      `json:"trace_id,omitempty"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
