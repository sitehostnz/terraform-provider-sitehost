package api

import (
	"context"
	"errors"
	"net/url"
	"time"
)

type JobService Client

func (c *Client) Job() *JobService { return (*JobService)(c) }

func (j *JobService) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(j).get(ctx, path, data, keys, v)
}

type JobType string

const (
	DaemonJob    JobType = "daemon"
	SchedulerJob JobType = "scheduler"
)

type Job struct {
	Created   time.Time `json:"created"`
	Started   time.Time `json:"created"`
	Completed time.Time `json:"created"`
	Message   string    `json:"message"`
	State     string    `json:"state"`
	Logs      []string  `json:"logs"`
}

// AddIP to a server.
func (j *JobService) Get(ctx context.Context, id uint, t JobType) (Job, error) {
	var r struct {
		Return  Job    `json:"return"`
		Message string `json:"msg"`
		Status  bool   `json:"status"`
	}
	err := j.get(ctx, "/job/get.json", url.Values{jobIDKey: []string{utoa(id)}, typeKey: []string{string(t)}}, []string{jobIDKey, typeKey}, &r)
	if err != nil {
		return Job{}, err
	}
	if !r.Status {
		return Job{}, errors.New(r.Message)
	}
	return r.Return, nil
}
