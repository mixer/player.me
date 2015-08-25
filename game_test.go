package player

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockFulfiller struct {
	mock.Mock
}

func (m *mockFulfiller) Request(method string, url string, params Query) ([]byte, error) {
	p, _ := params.toParams()
	args := m.Called(method, url, p.Encode())
	return args.Get(0).([]byte), args.Error(1)
}

func mockedClient() (*Client, *mockFulfiller) {
	f := new(mockFulfiller)
	c := New()
	c.Fulfiller = f

	return c, f
}

func TestGameList(t *testing.T) {
	client, f := mockedClient()
	f.On("Request", "GET", "https://player.me/api/v1/games", "").Return([]byte(`{
        "method": "GET",
        "pager": {
            "from": 0,
            "limit": 10,
            "total": 12345
        },
        "results": [
            {
                "alias": "Star Wars: Episode 1 - Battle for Naboo",
                "box": {
                    "cached": "//d1zqrvc06emslq.cloudfront.net/media/cache/game/ee/c7/3c/704b2c8ca7570dab157c6aecaeac4872.jpg",
                    "original": "//d1zqrvc06emslq.cloudfront.net/media/originals/game/ee/c7/3c/704b2c8ca7570dab157c6aecaeac4872.jpeg",
                    "original_filename": "704b2c8ca7570dab157c6aecaeac4872.jpeg"
                },
                "description": null,
                "favourites_count": 0,
                "has_favourited": false,
                "has_liked": false,
                "id": 617,
                "likes_count": 0,
                "slug": "star-wars-episode-i-battle-for-naboo",
                "title": "Star Wars: Episode I - Battle for Naboo",
                "url": "/g/star-wars-episode-i-battle-for-naboo"
            }
        ]
    }`), nil)

	result, err := client.GameList()
	assert.Nil(t, err)
	// Make sure it looks like it's right.
	assert.Equal(t, "GET", result.Method)
	assert.Equal(t, 12345, result.Pager.Total)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, "Star Wars: Episode 1 - Battle for Naboo", result.Results[0].Alias)
}

func TestGameListIterator(t *testing.T) {
	client, f := mockedClient()
	f.On("Request", "GET", "https://player.me/api/v1/games", "_from=0&_limit=100").Return([]byte(`{
        "method": "GET",
        "pager": {
            "from": 0,
            "limit": 100,
            "total": 300
        },
        "results": [{ "alias": "a" }]
    }`), nil)
	f.On("Request", "GET", "https://player.me/api/v1/games", "_from=100&_limit=100").Return([]byte(`{
        "method": "GET",
        "pager": {
            "from": 100,
            "limit": 100,
            "total": 300
        },
        "results": [{ "alias": "b" }]
    }`), nil)
	f.On("Request", "GET", "https://player.me/api/v1/games", "_from=200&_limit=100").Return([]byte(`{
        "method": "GET",
        "pager": {
            "from": 200,
            "limit": 100,
            "total": 300
        },
        "results": [{ "alias": "c" }]
    }`), nil)

	it := client.GameListIterate(false)
	out := []string{}
	for game := range it.Results {
		out = append(out, game.Alias)
	}
	it.Quit <- true
	assert.Equal(t, []string{"a", "b", "c"}, out)
}
