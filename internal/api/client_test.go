package api

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

var ctx = context.Background()

// transport returns mocked *http.Responses.
//
// Each element will be wrapped in a strings.Reader and set as
// http.Response.Body. (*transport)(nil), &(transport)(nil), and
// &(transport)([]string{}) will all panic to help signal an invalid test was
// configured.
type transport []string

// RoundTrip implements the http.RoundTripper interface for testing.
func (s *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	var body string
	body, *s = (*s)[0], (*s)[1:]
	return &http.Response{Body: ioutil.NopCloser(strings.NewReader(body))}, nil
}

// client returns a *Client with a transport that mocks the provided bodies.
func client(body ...string) *Client {
	return &Client{Client: http.Client{Transport: (*transport)(&body)}}
}

// ExampleClient shows basic usage for a client.
//
// This is an active integration test that will spawn and destroy production
// resources. As such, it is disabled when the -test.short flag is provided.
func ExampleClient() {
	if testing.Short() {
		return
	}
	id, err := atou(os.Getenv("SITEHOST_ID"))
	check(err)
	c := &Client{
		Base:   url.URL{Scheme: "https", Host: "mysth.safeserver.net.nz", Path: "1.0"},
		ID:     id,
		Token:  os.Getenv("SITEHOST_TOKEN"),
		Client: http.Client{},
	}
	_, err = c.Server().List(ctx)
	check(err)
	for i := range [5]struct{}{} {
		_, job, _, _, _, err := c.Server().Provision(ctx, "test", "AKLCITY", "XENLIT", "ubuntu-xenial.amd64", ParamIP{}, ParamName("example"+strconv.Itoa(i)))
		check(err)
		wait(c.Job(), job)
	}
	_, err = c.Server().List(ctx)
	check(err)

	// Output:
}

// check is a quick short circuit to fail if err != nil.
func check(err error) {
	if !errors.Is(err, nil) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// wait blocks until the job completes.
func wait(c *JobClient, job uint) {
	for {
		j, err := c.Get(ctx, job, DaemonJob)
		check(err)
		if j.State == CompleteJob {
			break
		}
	}
}
