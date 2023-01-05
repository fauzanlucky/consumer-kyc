package apicall

type Client struct {
	memberHost string
}

func NewClient(host string) *Client {
	return &Client{
		memberHost: host + "/member",
	}
}

type Caller interface {
}
