package player

const (
	// Request method constants
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"
)

type Client struct {
	// Base URL (defaults to https://player.me)
	Base string
	// The API fulfiller to use.
	Fulfiller Fulfiller
}

// Creates a new client.
func New() *Client {
	return &Client{
		Base:      "https://player.me",
		Fulfiller: &HttpFulfiller{},
	}
}

// Creates a new request instance.
func (c *Client) Request(method string, path string, params Query) *Request {
	return &Request{
		path:      path,
		method:    method,
		base:      c.Base,
		params:    params,
		fulfiller: c.Fulfiller,
	}
}
