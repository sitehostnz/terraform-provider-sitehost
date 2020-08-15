package api

import (
	"context"
	"errors"
	"net/url"
)

type APIService Client

func (c *Client) API() *APIService { return (*APIService)(c) }

func (s *APIService) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(s).get(ctx, path, data, keys, v)
}

// Info returns module access information about the API key in use.
func (s *APIService) Info(ctx context.Context) (modules []string, err error) {
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
