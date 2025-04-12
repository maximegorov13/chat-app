package tokencheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/internal/res"
)

type Client struct {
	serviceURL string
	httpClient *http.Client
}

type Config struct {
	ServiceURL string
}

func NewClient(conf Config) *Client {
	return &Client{
		serviceURL: conf.ServiceURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) IsTokenInvalid(ctx context.Context, token string) (*res.Response[auth.IsTokenInvalidResponse], error) {
	baseUrl := fmt.Sprintf("%s/api/auth/is-token-invalid", c.serviceURL)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiRes res.Response[auth.IsTokenInvalidResponse]
	if err = json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return nil, err
	}

	return &apiRes, nil
}
