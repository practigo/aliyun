package mts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/practigo/aliyun"
)

// API constants
const (
	Ver = "2014-06-18"
)

// api provides the common methods for a MTS API.
// It implements the aliyun.API interface.
type api struct {
	v url.Values
	aliyun.Base
}

func (api) Version() string {
	return Ver
}

func (a *api) Param() url.Values {
	return a.v
}

// SubmitJobsRequest contains the param for submitting a
// transcoding job. Only the OutputLocation is optional
// and default to "oss-cn-hangzhou".
type SubmitJobsRequest struct {
	Input          string `json:"Input"`
	OutputBucket   string `json:"OutputBucket"`
	OutputLocation string `json:"OutputLocation"`
	Outputs        string `json:"Outputs"`
	PipelineID     string `json:"PipelineId"`
}

// SubmitJobsAPI returns a API for SubmitJobs.
// https://help.aliyun.com/document_detail/29226.html
func SubmitJobsAPI(r *SubmitJobsRequest) aliyun.API {
	a := &api{v: url.Values{}}

	// api-specific mandotory params
	a.v.Add("Action", "SubmitJobs")
	a.v.Add("Input", r.Input)
	a.v.Add("OutputBucket", r.OutputBucket)
	a.v.Add("Outputs", r.Outputs)
	a.v.Add("PipelineId", r.PipelineID)

	// optional
	if r.OutputLocation != "" {
		a.v.Add("OutputLocation", r.OutputLocation)
	}

	return a
}

// QueryJobsAPI returns a API for QueryJobList.
// It accepts one or more (at most 10) JobIDs for query.
// https://help.aliyun.com/document_detail/29228.html
func QueryJobsAPI(id string, rest ...string) aliyun.API {
	a := &api{v: url.Values{}}

	// api-specific mandotory params
	a.v.Add("Action", "QueryJobList")
	ids := append(rest, id)
	a.v.Add("JobIds", strings.Join(ids, ","))

	fmt.Println(a.v.Encode())
	return a
}

// JobInfo represents the info for one job.
// The watermark & properties are too much so
// it's marshalled to raw bytes.
type JobInfo struct {
	JobID string `json:"JobId"`
	Input struct {
		Bucket   string `json:"Bucket"`
		Location string `json:"Location"`
		Object   string `json:"Object"`
	} `json:"Input"`
	Output struct {
		OutputFile struct {
			Bucket   string `json:"Bucket"`
			Location string `json:"Location"`
			Object   string `json:"Object"`
		} `json:"OutputFile"`
		TemplateID    string          `json:"TemplateId"`
		WaterMarkList json.RawMessage `json:"WaterMarkList,omitempty"` // it's too much
		Properties    json.RawMessage `json:"Properties,omitempty"`    // it's too much
		UserData      string          `json:"UserData"`
	} `json:"Output"`
	State        string    `json:"State"`
	Code         string    `json:"Code"`
	Message      string    `json:"Message"`
	Percent      int       `json:"Percent"`
	PipelineID   string    `json:"PipelineId"`
	CreationTime time.Time `json:"CreationTime"`
}

// JobResult gives the result of a job.
type JobResult struct {
	Success bool    `json:"Success"`
	Code    string  `json:"Code"`
	Message string  `json:"Message"`
	Job     JobInfo `json:"Job"`
}

// SubmitJobsResponse contains the response for SubmitJobs.
type SubmitJobsResponse struct {
	RequestID string `json:"RequestId"`
	List      struct {
		Result []JobResult `json:"JobResult"`
	} `json:"JobResultList"`
}

// QueryJobsResponse contains the response for QueryJobList.
type QueryJobsResponse struct {
	NonExistJobIDs struct {
		IDs []string `json:"String,omitempty"`
	} `json:"NonExistJobIds"`
	RequestID string `json:"RequestId"`
	JobList   struct {
		Job []JobInfo `json:"Job"`
	} `json:"JobList"`
}

// Submitter submits a transcoding job.
type Submitter interface {
	Submit(*SubmitJobsRequest) (SubmitJobsResponse, error)
}

// Querier queries a transcoding job.
type Querier interface {
	Query(id string, rest ...string) (QueryJobsResponse, error)
}

// Transcoder wraps a Submitter & a Querier.
type Transcoder interface {
	Submitter
	Querier
}

type transcoder struct {
	signer aliyun.Signer
	host   string
	cl     *http.Client
}

func (s *transcoder) Submit(r *SubmitJobsRequest) (resp SubmitJobsResponse, err error) {
	api := SubmitJobsAPI(r)
	err = aliyun.Get(s.cl, s.signer, api, s.host, &resp)
	return
}

func (s *transcoder) Query(id string, rest ...string) (resp QueryJobsResponse, err error) {
	api := QueryJobsAPI(id, rest...)
	err = aliyun.Get(s.cl, s.signer, api, s.host, &resp)
	return
}

// New returns a new Transcoder.
func New(s aliyun.Signer, host string) Transcoder {
	return &transcoder{
		signer: s,
		host:   host,
		cl: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
