package model

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type DataResponse struct {
	Data string `json:"data"`
}

func NewSuccessResponse() *SuccessResponse {
	resp := SuccessResponse{Message: "success"}
	return &resp
}
