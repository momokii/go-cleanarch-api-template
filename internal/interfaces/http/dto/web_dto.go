package dto

type PaginationData struct {
	TotalData int
	Data      interface{}
}

type WebResponse struct {
	Errors  bool        `json:"errors"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
