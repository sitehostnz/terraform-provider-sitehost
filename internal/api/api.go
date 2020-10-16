package api

import (
	"context"
	"errors"
	"net/url"
)

// APIClient is the container with methods to call /api/.
//
//nolint:golint
type APIClient Client

// API vends a client for making calls to /api/.
func (c *Client) API() *APIClient { return (*APIClient)(c) }

// get is a wrapper to prevent inline type conversions.
//
// See: Client.get.
func (s *APIClient) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(s).get(ctx, path, data, keys, v)
}

// Info returns module access information about the API key in use.
func (s *APIClient) Info(ctx context.Context) (modules []string, err error) {
	var r struct {
		Return struct {
			Modules []string `json:"modules"`
		} `json:"return"`
		Status  bool   `json:"status"`
		Message string `json:"msg"`
	}
	err = s.get(ctx, "/api/get_info.json", nil, nil, &r)
	if err != nil {
		return nil, err
	}
	if !r.Status {
		return nil, errors.New(r.Message)
	}
	return r.Return.Modules, nil
}
