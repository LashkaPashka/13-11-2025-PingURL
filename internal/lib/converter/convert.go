package converter

import (
	"github.com/LashkaPashka/LinkCheck/internal/http-server/handlers/UrlStatus/payload"
	"github.com/LashkaPashka/LinkCheck/internal/model"
	"github.com/google/uuid"
)

const (
	statusPending = "pending"
)

func Convert(payload payload.Request) model.Links {
	var url []model.Url

	for _, link := range payload.Links {
		url = append(url, model.Url{
			UrlName: link,
			Status: statusPending,
		})
	}

	return model.Links{
		Urls: url,
		Status: statusPending,
		NumberLink: uuid.New().String(),
	}
}