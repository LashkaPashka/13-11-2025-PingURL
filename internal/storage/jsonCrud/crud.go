package jsoncrud

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	"github.com/LashkaPashka/LinkCheck/internal/lib/encode"
	"github.com/LashkaPashka/LinkCheck/internal/model"
)

type Storage struct {
	storagePath string
	logger      *slog.Logger
}

func New(storagePath string, logger *slog.Logger) *Storage {
	if _, err := os.Stat(storagePath); err != nil {
		if os.IsNotExist(err) {
			logger.Error("Storage path does not exist",
				slog.String("path", storagePath),
			)
			return nil
		}

		logger.Error("Error checking storage path",
			slog.String("path", storagePath),
			slog.String("err", err.Error()),
		)
		return nil
	}

	return &Storage{
		storagePath: storagePath,
		logger:      logger,
	}
}

func (s *Storage) Save(link model.Links) (success bool, err error) {
	const op = "LinkCheck.storage.crud.Save"

	file, err := os.ReadFile(s.storagePath)
	if err != nil {
		s.logger.Error("Failed to read file",
			slog.String("op", op),
			slog.String("err", err.Error()))

		return false, err
	}

	links := []model.Links{}

	json.Unmarshal(file, &links)

	links = append(links, link)

	encLinks, err := encode.Encode(links, s.logger)
	if err != nil {
		return false, err
	}

	if err = os.WriteFile(s.storagePath, encLinks, 0644); err != nil {
		s.logger.Error("Invalid save task in file",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return false, err
	}

	return true, nil
}

func (s *Storage) Update(links []model.Links) (bool, error) {
    const op = "LinkCheck.storage.crud.Update"

    data, err := json.MarshalIndent(links, "", "  ")
    if err != nil {
        return false, err
    }

    tmpFile := s.storagePath + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return false, err
    }
    if err := os.Rename(tmpFile, s.storagePath); err != nil {
        return false, err
    }

    s.logger.Info("Link updated successfully",
        slog.String("op", op),
    )
    return true, nil
}

func (s *Storage) ReadByNumID(targetLinkNum string) (link model.Links, err error) {
	links, err := s.readAllLinks()
	if err != nil {
		return link, err
	}

	for _, l := range links {
		if strings.Compare(l.NumberLink, targetLinkNum) == 0 {
			link = l
			break
		}
	}

	return link, err
}

func (s *Storage) ReadAll() (links []model.Links, err error) {
	return s.readAllLinks()
}

func (s *Storage) readAllLinks() ([]model.Links, error) {
	file, err := os.Open(s.storagePath)
	if err != nil {
		s.logger.Error("Failed to open storage file",
			slog.String("file", s.storagePath),
			slog.String("err", err.Error()),
		)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var b []byte
	for scanner.Scan() {
		b = append(b, []byte(scanner.Text()+"\n")...)
	}

	if err := scanner.Err(); err != nil {
		s.logger.Error("Scanner encountered an error",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	var links []model.Links
	json.Unmarshal(b, &links)

	return links, nil
}