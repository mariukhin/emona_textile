package api

import "math"

type PageResponse struct {
	Pages  int         `json:"pages"`
	Total  int         `json:"total"`
	Result interface{} `json:"result"`
}

func NewPageResp(total, pageSize int, result interface{}) PageResponse {
	return PageResponse{
		Pages:  int(math.Ceil(float64(total) / float64(pageSize))),
		Total:  total,
		Result: result,
	}
}
