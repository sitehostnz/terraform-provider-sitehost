package api

type voidResponse struct {
	Message string `json:"msg"`
	Status  bool   `json:"status"`
}

type jobIDResponse struct {
	Return struct {
		JobID uint `json:"job_id,string"`
	} `json:"return"`
	Message string `json:"msg"`
	Status  bool   `json:"status"`
}
