package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"

	generatepdf "github.com/LashkaPashka/LinkCheck/internal/lib/generatePDF"
	statusping "github.com/LashkaPashka/LinkCheck/internal/lib/statusPing"
	"github.com/LashkaPashka/LinkCheck/internal/model"
)

const (
	statusFailed  = "failed"
	statusPending = "pending"
	statusDone    = "done"
	NotFound      = 404
	Up            = "up"
	Down          = "down"
)

type Storager interface {
	Save(link model.Links) (success bool, err error)
	ReadByNumID(targetLinkNum string) (link model.Links, err error)
	ReadAll() (link []model.Links, err error)
	Update(links []model.Links) (success bool, err error)
}

type Service struct {
	logger  *slog.Logger
	storage Storager
}

func New(storage Storager, logger *slog.Logger) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
	}
}

func (s *Service) Save(link model.Links) (updatedLink model.Links, success bool, err error) {
	const op = "LinkCheck.service.Save"

	success, err = s.storage.Save(link)
	if err != nil {
		s.logger.Error("Failed to save link",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return model.Links{}, false, err
	}

	if err := s.FoundPendingLink(); err != nil {
		s.logger.Error("Failed to find pending link",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return model.Links{}, false, err
	}

	updatedLink, err = s.storage.ReadByNumID(link.NumberLink)
	if err != nil {
		s.logger.Error("Failed to read targetlink",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return model.Links{}, false, err
	}

	return updatedLink, success, err
}

func (s *Service) pingUrl(link model.Links) (model.Links, error) {
	const op = "LinkCheck.service.pingUrl"

	for i := 0; i < len(link.Urls); i++ {
		if strings.Compare(link.Urls[i].Status, statusPending) != 0 {
			continue
		}

		statusCode, err := statusping.Ping(link.Urls[i].UrlName, s.logger)
		if err != nil {
			var dnsError *net.DNSError
			if errors.As(err, &dnsError) && dnsError.IsNotFound {
				link.Urls[i].Available = Down
				link.Urls[i].Status = statusFailed
				link.Urls[i].StatusCode = NotFound
				continue
			}

			s.logger.Error("Failed to ping URL",
				slog.String("op", op),
				slog.String("url", link.Urls[i].UrlName),
				slog.String("err", err.Error()),
			)

			return model.Links{}, err
		}

		link.Urls[i].Available = Up
		link.Urls[i].Status = statusDone
		link.Urls[i].StatusCode = statusCode
	}

	link.Status = statusDone

	return link, nil
}

func (s *Service) CreatePDF(listLinks ...string) (buf []byte, err error) {
	const op = "LinkCheck.service.CreatePDF"

	var urls []model.Url

	for _, numLinkID := range listLinks {
		link, err := s.storage.ReadByNumID(numLinkID)
		if err != nil {
			s.logger.Error("Failed to read link from storage",
				slog.String("op", op),
				slog.String("err", err.Error()),
			)

			return buf, err
		}

		if len(link.Urls) == 0 && link.NumberLink == "" {
			continue
		}

		urls = append(urls, link.Urls...)
	}

	buf, err = generatepdf.Generate(urls)
	if err != nil {
		s.logger.Error("Failed to generate PDF document",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return buf, err
	}

	return buf, err
}

func (s *Service) FoundPendingLink() (err error) {
	const op = "LinkCheck.service.FoundPendingTask"

	links, err := s.storage.ReadAll()
	if err != nil {
		s.logger.Error("Failed to read all links from storage",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return err
	}

	if links == nil {
		return fmt.Errorf("storage is empty; links: %v", links)
	}

	for idx := range links {
		if links[idx].Status == statusDone {
			continue
		}

		newLink, err := s.pingUrl(links[idx])
		if err != nil {
			s.logger.Error("Failed to ping link URL",
				slog.String("op", op),
				slog.String("err", err.Error()),
			)

			return err
		}

		links[idx] = newLink
	}

	success, err := s.storage.Update(links)
	if err != nil || !success {
		s.logger.Error("Failed to update link",
			slog.String("op", op),
			slog.Bool("success", success),
			slog.String("err", err.Error()),
		)

		return err
	}

	return err
}
