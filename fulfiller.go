package player

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// The fulfiller is used to run requests to the player.me API. You
// can implement it and pass it into the client if you want to
// override some logic.
type Fulfiller interface {
	Request(method string, url string, params Query) ([]byte, error)
}

// Fulfiller for the HTTP API.
type HttpFulfiller struct{}

func (h HttpFulfiller) resolveRequest(method string, url string, params url.Values) (*http.Request, error) {
	if method == GET {
		return http.NewRequest(method, url+"?"+params.Encode(), nil)
	}

	encoded := params.Encode()
	r, err := http.NewRequest(method, url, bytes.NewBufferString(encoded))
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(encoded)))

	return r, nil
}

// Runs a request against the player.me REST API.
func (h HttpFulfiller) Request(method string, target string, params Query) ([]byte, error) {
	var values url.Values
	if params != nil {
		v, err := params.toParams()
		if err != nil {
			return nil, err
		}
		values = v
	} else {
		values = url.Values{}
	}

	req, err := h.resolveRequest(method, target, values)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return out, err
}
