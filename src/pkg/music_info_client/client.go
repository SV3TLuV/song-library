package music_info_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"song-library-api/src/pkg/music_info_client/model"
)

type MusicInfoClient struct {
	baseUrl    string
	httpClient *http.Client
}

func NewMusicInfoClient(baseUrl string) *MusicInfoClient {
	return &MusicInfoClient{
		baseUrl:    baseUrl,
		httpClient: &http.Client{},
	}
}

func (c *MusicInfoClient) GetSongInfo(ctx context.Context, group, song string) (*model.SongDetail, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/info", c.baseUrl))
	if err != nil {
		return nil, errors.Wrap(err, "url parse failed")
	}

	query := requestURL.Query()
	query.Set("group", group)
	query.Set("song", song)
	requestURL.RawQuery = query.Encode()

	response, err := c.httpClient.Get(requestURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("request failed with status code %d", response.StatusCode)
	}

	var songDetail model.SongDetail
	if err := json.NewDecoder(response.Body).Decode(&songDetail); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &songDetail, nil
}
