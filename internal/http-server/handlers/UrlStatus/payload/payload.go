package payload

import "github.com/LashkaPashka/LinkCheck/internal/model"

type Request struct {
	Links []string `json:"links"`
}

type Response struct {
	Links      []model.Url `json:"links"`
	Status     string      `json:"status"`
	NumberLink string      `json:"links_num"`
}
