package photo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ras0q/goalie"
)

const photoAPIEndpoint = "https://jsonplaceholder.typicode.com/photos"

type Photo struct {
	AlbumID      int    `json:"albumId"`
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type Service struct{}

func NewPhotoService() *Service {
	return &Service{}
}

func (*Service) GetPhoto(id int) (_ *Photo, err error) {
	g := goalie.New()
	defer g.Collect(&err)

	resp, err := http.Get(fmt.Sprintf("%s/%d", photoAPIEndpoint, id))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer g.Guard(resp.Body.Close)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var photo Photo
	err = json.Unmarshal(body, &photo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &photo, nil
}
