package urlstatus

import (
	"log/slog"
	"net/http"

	"github.com/LashkaPashka/LinkCheck/internal/http-server/handlers/UrlStatus/payload"
	"github.com/LashkaPashka/LinkCheck/internal/lib/converter"
	"github.com/LashkaPashka/LinkCheck/internal/lib/req"
	"github.com/LashkaPashka/LinkCheck/internal/lib/res"
	"github.com/LashkaPashka/LinkCheck/internal/model"
)

type Service interface {
	Save(link model.Links) (updatedLink model.Links, success bool, err error)
}

func New(service Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.Request](w, r, logger)
		if err != nil {
			logger.Error("Invalid validation",
					slog.Any("body", body),
					slog.String("err", err.Error()),
				)
			res.Encode(w, body, http.StatusUnprocessableEntity)
			return
		}

		links := converter.Convert(body)

		link, success, err := service.Save(links)
		if !success || err != nil {
			logger.Error("Failed to save links",
					slog.Any("body", body),
					slog.String("err", err.Error()),
				)
			res.Encode(w, body, http.StatusConflict)
			return
		}

		res.Encode(w, &payload.Response{
			Links: link.Urls,
			Status: link.Status,
			NumberLink: link.NumberLink,
		}, http.StatusOK)
	}
}