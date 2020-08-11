// Package api implements a client for the sitehost mysth.safeserver.net.nz service.
//
// Consumers should begin with a Client, then use converter methods to access
// the api of their choice:
//	m, er := (&Client{}).API().Info(ctx)
//
// Converters are just syntax sugar for type conversions:
//	c := &Client{}
//	c.API() == APIClient(c)
//
// Because all Clients share the same memory layout you can also start directly
// with the desired client:
//	j, err := (&JobClient{}).Get(ctx, 4, DaemonJob)
//
// Currently Client.Token is required for all calls, and Client.ID for most. It
// is suggested not to bake tokens into source, instead prefering environment
// variables. At a later date Client.Token == "" may extract a default
// environment variable.
//	&Client{Token: os.Getenv("SITEHOST_TOKEN"}, ID: 1234}
//
// For development and custom environments it may be required to specify a
// custom base url to override Base or a custom http.Client. E.g. during
// testing we override Client.Client.Transport to mock api calls. These
// optional parameters can be specified like so:
//	&Client{
//		Base: url.URL{Scheme: "http", Host: "dev.safeserver.net.nz"},
//		Client: http.Client{Transport:
//			http.NewFileTransport(http.Dir("/path/to/mocks"))
//		},
//	}
//
// While this interface is currently /internal/, it can be expanded into a full
// published api client at a later date.
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

// Base is the url prefix of the service encoded as a string.
const Base = "https://mysth.safeserver.net.nz/1.0"

// encode is a modified url.Values.Encode that supports custom value ordering.
//
// Client level informaiton, Token and ID, is also added where appropriate.
//
// keys are accepted as input, as opposed to url.Values, which generates a
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

// Client is interacts with the sitehost api.
type Client struct {
	// Base is the url prefix of the api service.
	//
	// See: Base, it is both the default and an example.
	Base url.URL

	// Token is the apikey for your client.
	Token string

	// ID is the client_id you are making calls about.
	//
	// This may not be sent for all api calls, since it is not always required.
	ID uint

	// Client can be used to implement or override default routing behaviour.
	Client http.Client
}

// base conjoins the provided path, with the the calculated base url.
//
// If Client.Base is unconfigured, Base is used.
func (c Client) base(path string) string {
	base := Base
	if c.Base != (url.URL{}) {
		base = c.Base.String()
	}
	return base + path
}

// unmarshal consumes an *http.Response and encodes the json body into v.
//
// In testing, all responses from the service are http.StatusOK (200), and in
// json format.  For ease of use, it also consumes an err, returning it if
// non-nil.
func unmarshal(r *http.Response, e error, v interface{}) (err error) {
	if e != nil {
		return e
	}
	defer func() {
		e = r.Body.Close()
		if err == nil {
			err = e
		}
	}()
	return json.NewDecoder(r.Body).Decode(v)
}

// utoa is equivalent to strconv.FormatUint(uint64(u), 10).
//
// It is based on strconv.Itoa.
func utoa(u uint) string { return strconv.FormatUint(uint64(u), 10) }

// postForm is a custom implementation of ctxhttp.PostForm that allows for custom parameter ordering.
func (c *Client) postForm(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	r, err := ctxhttp.Post(ctx, &c.Client, c.base(path), "application/x-www-form-urlencoded", strings.NewReader(c.encode(data, keys)))
	return unmarshal(r, err, v)
}

// get is a custom ctxhttp.Get that encodes parameters in custom order.
func (c *Client) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	r, err := ctxhttp.Get(ctx, &c.Client, c.base(path)+"?"+c.encode(data, keys))
	return unmarshal(r, err, v)
}
