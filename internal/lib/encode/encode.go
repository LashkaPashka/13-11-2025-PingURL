package encode

import (
	"encoding/json"
	"log/slog"
)

func Encode(payload any, logger *slog.Logger) ([]byte, error) {
	encPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Invalid encoding payload",
			slog.Any("payload", payload),
			slog.String("err", err.Error()),
	    )

		return nil, err
	}

	return encPayload, nil
}