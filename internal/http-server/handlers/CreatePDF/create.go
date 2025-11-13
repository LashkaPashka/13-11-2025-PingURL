package createpdf

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/LashkaPashka/LinkCheck/internal/http-server/handlers/CreatePDF/payload"
	"github.com/LashkaPashka/LinkCheck/internal/lib/req"
	"github.com/LashkaPashka/LinkCheck/internal/lib/res"
)

type Service interface {
	CreatePDF(listLinks ...string) (buf []byte, err error)
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

		buf, err := service.CreatePDF(body.LinksList...)
		if err != nil {
			logger.Error("Failed to save links",
					slog.Any("body", body),
					slog.String("err", err.Error()),
				)
			res.Encode(w, body, http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(buf)))

		w.Write(buf)
	}
}