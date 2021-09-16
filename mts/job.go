package mts

import (
	"encoding/json"
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

// An api provides the common parts for a MTS API.
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

// A SubmitJobsRequest contains the param for submitting a
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

	// fmt.Println(a.v.Encode())
	return a
}

// JobIO is the job input/output.
type JobIO struct {
	Bucket   string `json:"Bucket"`
	Location string `json:"Location"`
	Object   string `json:"Object"`
}

// A JobInfo represents the info for one job.
// The Output has too many fields so it's marshalled to raw bytes.
type JobInfo struct {
	JobID        string          `json:"JobId"`
	Input        JobIO           `json:"Input"`
	Output       json.RawMessage `json:"Output"`
	State        string          `json:"State"`
	Code         string          `json:"Code"`
	Message      string          `json:"Message"`
	Percent      int             `json:"Percent"`
	PipelineID   string          `json:"PipelineId"`
	CreationTime time.Time       `json:"CreationTime"`
	FinishTime   time.Time       `json:"FinishTime"`
}

// JobOutputInfo is a mapping for JobInfo.Output.
// This is NOT mean to be completed.
type JobOutputInfo struct {
	OutputFile JobIO           `json:"OutputFile"`
	UserData   string          `json:"UserData"`
	Priority   string          `json:"Priority"`
	Properties json.RawMessage `json:"Properties"`
	ExtendData string          `json:"ExtendData"`
	TemplateID string          `json:"TemplateId"`
}

// A JobResult gives the result of a job.
type JobResult struct {
	Success bool    `json:"Success"`
	Code    string  `json:"Code"`
	Message string  `json:"Message"`
	Job     JobInfo `json:"Job"`
}

// A SubmitJobsResponse contains the response for SubmitJobs.
type SubmitJobsResponse struct {
	RequestID string `json:"RequestId"`
	List      struct {
		Result []JobResult `json:"JobResult"`
	} `json:"JobResultList"`
}

// A QueryJobsResponse contains the response for QueryJobList.
type QueryJobsResponse struct {
	NonExistJobIDs struct {
		IDs []string `json:"String,omitempty"`
	} `json:"NonExistJobIds"`
	RequestID string `json:"RequestId"`
	JobList   struct {
		Job []JobInfo `json:"Job"`
	} `json:"JobList"`
}

// A Submitter submits a transcoding job.
type Submitter interface {
	Submit(*SubmitJobsRequest) (SubmitJobsResponse, error)
}

// A Querier queries a transcoding job.
type Querier interface {
	Query(id string, rest ...string) (QueryJobsResponse, error)
}

// A Transcoder wraps a Submitter & a Querier.
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

// New returns a new Transcoder with a 10s-timeout
// HTTP client.
func New(s aliyun.Signer, host string) Transcoder {
	return &transcoder{
		signer: s,
		host:   host,
		cl:     aliyun.TimeoutClient(10 * time.Second),
	}
}
