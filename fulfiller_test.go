package player

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestWorksWithBlankGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "", r.URL.RawQuery)
	}))
	defer ts.Close()

	HttpFulfiller{}.Request("GET", ts.URL, nil)
}

func TestWorksWithGetParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "_limit=100", r.URL.RawQuery)
	}))
	defer ts.Close()

	HttpFulfiller{}.Request("GET", ts.URL, Query{"_limit": 100})
}

func TestWorksWithBlankPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, url.Values{}, r.Form)
	}))
	defer ts.Close()

	HttpFulfiller{}.Request("POST", ts.URL, nil)
}

func TestWorksWithPostParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, url.Values{"_limit": []string{"100"}}, r.Form)
	}))
	defer ts.Close()

	HttpFulfiller{}.Request("POST", ts.URL, Query{"_limit": 100})
}
