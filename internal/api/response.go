package api

// voidResponse is the syntax of apis that do not pass back any data.
type voidResponse struct {
	Message string `json:"msg"`
	Status  bool   `json:"status"`
}

// jobIDResponse is a standard syntax for async api calls that return a job id.
//
// See: JobClient.Get for polling the status of a job.
type jobIDResponse struct {
	Return struct {
		JobID uint `json:"job_id,string"`
	} `json:"return"`
	Message string `json:"msg"`
	Status  bool   `json:"status"`
}
