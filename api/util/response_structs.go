package util

type HttpResponse struct {
	Items   interface{} `json:"items"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}
