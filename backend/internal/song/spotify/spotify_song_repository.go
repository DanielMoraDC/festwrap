package spotify

import (
	"fmt"
	"net/url"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
	"festwrap/internal/song/errors"
)

type SpotifySongRepository struct {
	accessToken  string
	host         string
	httpSender   httpsender.HTTPRequestSender
	deserializer serialization.Deserializer[[]song.Song]
}

func NewSpotifySongRepository(
	accessToken string,
	httpSender httpsender.HTTPRequestSender,
) *SpotifySongRepository {
	return &SpotifySongRepository{
		accessToken: accessToken, host: "api.spotify.com", httpSender: httpSender, deserializer: &SpotifySongsDeserializer{}}
}

func (r *SpotifySongRepository) GetSong(artist string, title string) (*song.Song, error) {
	httpOptions := r.createSongHttpOptions(artist, title)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return nil, errors.NewCannotRetrieveSongError(err.Error())
	}

	songs, err := r.deserializer.Deserialize(*responseBody)
	if err != nil {
		return nil, errors.NewCannotRetrieveSongError(err.Error())
	}

	allSongs := *songs
	if len(allSongs) == 0 {
		errorMsg := fmt.Sprintf("No songs found for song %s (%s)", title, artist)
		return nil, errors.NewCannotRetrieveSongError(errorMsg)
	}

	// We assume the first result is the most trusted one
	result := allSongs[0]
	return &result, nil
}

func (r *SpotifySongRepository) SetDeserializer(deserializer serialization.Deserializer[[]song.Song]) {
	r.deserializer = deserializer
}

func (r *SpotifySongRepository) createSongHttpOptions(artist string, title string) httpsender.HTTPRequestOptions {
	httpOptions := httpsender.NewHTTPRequestOptions(r.getSetlistFullUrl(artist, title), httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", r.accessToken)},
	)
	return httpOptions
}

func (r *SpotifySongRepository) getSetlistFullUrl(artist string, title string) string {
	queryParams := url.Values{}
	queryParams.Set("q", fmt.Sprintf("artist:%s track:%s", artist, title))
	queryParams.Set("type", "track")
	setlistPath := "v1/search"
	return fmt.Sprintf("https://%s/%s?%s", r.host, setlistPath, queryParams.Encode())
}
