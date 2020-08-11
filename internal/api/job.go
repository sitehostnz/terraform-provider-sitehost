package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"time"
)

// JobClient is the container with methods to call /job/.
type JobClient Client

// Job vends a client for making calls to /job/.
func (c *Client) Job() *JobClient { return (*JobClient)(c) }

// get is a wrapper to prevent inline type conversions.
//
// See: Client.get.
func (j *JobClient) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(j).get(ctx, path, data, keys, v)
}

// JobType is a category of job.
type JobType string

// Enumerated JobTypes.
const (
	DaemonJob    JobType = "daemon"
	SchedulerJob JobType = "scheduler"
)

// JobState is the state of a given job.
type JobState string

// Common JobStates.
const (
	PendingJob  JobState = "Pending"
	CompleteJob JobState = "Completed"
)

// Job is the metadata about a job.
type Job struct {
	Created   time.Time
	Started   time.Time
	Completed time.Time
	Message   string
	State     JobState
	Logs      []struct {
		Date    time.Time
		Level   uint
		Message string
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Job) UnmarshalJSON(b []byte) error {
	var v struct {
		Created   *stime    `json:"created"`
		Started   *stime    `json:"started"`
		Completed *stime    `json:"completed"`
		Message   *string   `json:"message"`
		State     *JobState `json:"state"`
		Logs      []struct {
			Date    stime  `json:"date"`
			Level   uint   `json:"level,string"`
			Message string `json:"message"`
		} `json:"logs"`
	}
	v.Created = (*stime)(&j.Created)
	v.Started = (*stime)(&j.Started)
	v.Completed = (*stime)(&j.Completed)
	v.Message = &j.Message
	v.State = &j.State
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	for _, l := range v.Logs {
		j.Logs = append(j.Logs, struct {
			Date    time.Time
			Level   uint
			Message string
		}{time.Time(l.Date), l.Level, l.Message})
	}
	return nil
}

// Get retrieve information about a given job id.
func (j *JobClient) Get(ctx context.Context, id uint, t JobType) (Job, error) {
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
