package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context/ctxhttp"
)

var Base = "https://mysth.safeserver.net.nz/1.0"

// encode is a modified url.Values.Encode that supports custom value ordering.
//
// Client level informaiton, Token and ID, is also added where appropriate.
//
// keys is accepted as input, as opposed to url.Values, which generates a
// sorted list. Extra or missing keys be skipped.
func (c Client) encode(v url.Values, keys []string) string {
	if v == nil {
		v = url.Values{}
	}
	v.Set(apiKey, c.Token)
	v.Set(clientIDKey, utoa(c.ID))
	var buf strings.Builder
	for _, k := range append([]string{apiKey}, keys...) {
		vs, ok := v[k]
		if !ok {
			continue
		}
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

type Client struct {
	Base  string
	Token string
	ID    uint

	Client http.Client
}

func (c Client) base(path string) string {
	if c.Base != "" {
		return c.Base + path
	}
	return Base + path
}

func unmarshal(r *http.Response, err error, v interface{}) error {
	if err != nil {
		return err
	}
	defer func() { _ = r.Body.Close() }()
	return json.NewDecoder(r.Body).Decode(v)
}

func utoa(u uint) string { return strconv.FormatUint(uint64(u), 10) }

func (c *Client) postForm(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	r, err := ctxhttp.Post(ctx, &c.Client, c.base(path), "application/x-www-form-urlencoded", strings.NewReader(c.encode(data, keys)))
	return unmarshal(r, err, v)
}

func (c *Client) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	r, err := ctxhttp.Get(ctx, &c.Client, c.base(path)+"?"+c.encode(data, keys))
	return unmarshal(r, err, v)
}
