package model

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Ttl   int64  `json:"ttl"`
}

type DataResponse struct {
	Data string `json:"data"`
}
