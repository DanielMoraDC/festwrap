package spotify

import (
	"errors"
	"testing"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/playlist"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
)

func fakeSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	emptyResponse := []byte("")
	sender.SetResponse(&emptyResponse)
	return &sender
}

func errorSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test send error"))
	return &sender
}

func defaultPlaylistId() string {
	return "test_id"
}

func defaultSongs() []song.Song {
	return []song.Song{song.NewSong("uri1"), song.NewSong("uri2")}
}

func defaultUserId() string {
	return "some_user_id"
}

func defaultPlaylist() playlist.Playlist {
	return playlist.Playlist{Name: "my-playlist", Description: "some playlist", IsPublic: false}
}

func defaultSongsBody() []byte {
	return []byte(`{"uris":["uri1","uri2"]}`)
}

func defaultPlaylistBody() []byte {
	return []byte(`{"name":"my-playlist","description":"some playlist","is_public":false}`)
}

func defaultSongsSerializer() *serialization.FakeSerializer[SpotifySongs] {
	serializer := serialization.FakeSerializer[SpotifySongs]{}
	serializer.SetResponse(defaultSongsBody())
	return &serializer
}

func errorSongsSerializer() *serialization.FakeSerializer[SpotifySongs] {
	serializer := defaultSongsSerializer()
	serializer.SetError(errors.New("test songs error"))
	return serializer
}

func defaultPlaylistSerializer() *serialization.FakeSerializer[SpotifyPlaylist] {
	serializer := serialization.FakeSerializer[SpotifyPlaylist]{}
	serializer.SetResponse(defaultPlaylistBody())
	return &serializer
}

func errorPlaylistSerializer() *serialization.FakeSerializer[SpotifyPlaylist] {
	serializer := defaultPlaylistSerializer()
	serializer.SetError(errors.New("test playlist error"))
	return serializer
}

func expectedAddSongsHttpOptions() httpsender.HTTPRequestOptions {
	options := httpsender.NewHTTPRequestOptions("https://api.spotify.com/v1/playlists/test_id/tracks", httpsender.POST, 201)
	options.SetHeaders(defaultHeaders())
	options.SetBody(defaultSongsBody())
	return options
}

func expectedCreatePlaylistHttpOptions() httpsender.HTTPRequestOptions {
	options := httpsender.NewHTTPRequestOptions("https://api.spotify.com/v1/users/some_user_id/playlists", httpsender.POST, 201)
	options.SetHeaders(defaultHeaders())
	options.SetBody(defaultPlaylistBody())
	return options
}

func defaultHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer abcdefg12345",
		"Content-Type":  "application/json",
	}
}

func spotifyPlaylistRepository() SpotifyPlaylistRepository {
	repository := NewSpotifyPlaylistRepository(fakeSender(), "abcdefg12345")
	repository.SetSongSerializer(defaultSongsSerializer())
	repository.SetPlaylistSerializer(defaultPlaylistSerializer())
	return repository
}

func TestAddSongsReturnsErrorWhenNoSongsProvided(t *testing.T) {
	repository := spotifyPlaylistRepository()

	err := repository.AddSongs(defaultPlaylistId(), []song.Song{})

	testtools.AssertErrorIsNotNil(t, err)
}

func TestAddSongsSerializesInputSongs(t *testing.T) {
	repository := spotifyPlaylistRepository()
	serializer := defaultSongsSerializer()
	repository.SetSongSerializer(serializer)

	repository.AddSongs(defaultPlaylistId(), defaultSongs())

	expected := SpotifySongs{Uris: []string{"uri1", "uri2"}}
	actual := serializer.GetArgs()
	testtools.AssertEqual(t, actual, expected)
}

func TestAddSongsReturnsErrorOnNonSerializationError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	repository.SetSongSerializer(errorSongsSerializer())

	err := repository.AddSongs(defaultPlaylistId(), defaultSongs())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestAddSongsSendsRequestUsingProperOptions(t *testing.T) {
	sender := fakeSender()
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(sender)

	repository.AddSongs(defaultPlaylistId(), defaultSongs())

	actual := sender.GetSendArgs()
	testtools.AssertEqual(t, actual, expectedAddSongsHttpOptions())
}

func TestAddSongsReturnsErrorOnSendError(t *testing.T) {
	sender := errorSender()
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(sender)

	err := repository.AddSongs(defaultPlaylistId(), defaultSongs())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestAddSongsSerializesInputPlaylist(t *testing.T) {
	repository := spotifyPlaylistRepository()
	serializer := defaultPlaylistSerializer()
	repository.SetPlaylistSerializer(serializer)

	repository.CreatePlaylist(defaultUserId(), defaultPlaylist())

	expected := SpotifyPlaylist{Name: "my-playlist", Description: "some playlist", IsPublic: false}
	actual := serializer.GetArgs()
	testtools.AssertEqual(t, actual, expected)
}

func TestCreatePlaylistReturnsErrorOnPlaylistSerializationError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	repository.SetPlaylistSerializer(errorPlaylistSerializer())

	err := repository.CreatePlaylist(defaultUserId(), defaultPlaylist())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestCreatePlaylistSendsCreateRequestWithOptions(t *testing.T) {
	sender := fakeSender()
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(sender)

	repository.CreatePlaylist(defaultUserId(), defaultPlaylist())

	actual := sender.GetSendArgs()
	testtools.AssertEqual(t, actual, expectedCreatePlaylistHttpOptions())
}

func TestCreatePlaylistReturnsErrorOnSendError(t *testing.T) {
	sender := errorSender()
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(sender)

	err := repository.CreatePlaylist(defaultUserId(), defaultPlaylist())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestAddSongsPlaylistSendsOptionsUsingSerializerIntegration(t *testing.T) {
	sender := fakeSender()
	serializer := serialization.NewJsonSerializer[SpotifySongs]()
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(sender)
	repository.SetSongSerializer(&serializer)

	repository.AddSongs(defaultPlaylistId(), defaultSongs())

	actual := sender.GetSendArgs()
	testtools.AssertEqual(t, actual, expectedAddSongsHttpOptions())
}

func TestCreatePlaylistSendsOptionsUsingSerializerIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	sender := fakeSender()
	serializer := serialization.NewJsonSerializer[SpotifyPlaylist]()
	repository := spotifyPlaylistRepository()
	repository.SetPlaylistSerializer(&serializer)
	repository.SetHTTPSender(sender)

	repository.CreatePlaylist(defaultUserId(), defaultPlaylist())

	actual := sender.GetSendArgs()
	testtools.AssertEqual(t, actual, expectedCreatePlaylistHttpOptions())
}
