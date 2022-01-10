package helpers

type ApiResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewApiResponse(status int, message string, data interface{}) *ApiResponse {
	return &ApiResponse{Status: status, Message: message, Data: data}
}
