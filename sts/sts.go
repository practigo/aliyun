package sts

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/practigo/aliyun"
)

// API constants
const (
	Host = "https://sts.aliyuncs.com/" // public domain
	Ver  = "2015-04-01"
)

// A Credentials is the credentials obtained by AssumedRole.
type Credentials struct {
	AccessKeySecret string    `json:"AccessKeySecret"`
	AccessKeyID     string    `json:"AccessKeyId"`
	Expiration      time.Time `json:"Expiration"` // "2006-01-02T15:04:05Z"
	SecurityToken   string    `json:"SecurityToken"`
}

// An AssumedRoleUser is the user representation.
type AssumedRoleUser struct {
	AssumedRoleID string `json:"AssumedRoleId"`
	Arn           string `json:"Arn"`
}

// An AssumeRoleResponse is the response for action AssumeRole.
type AssumeRoleResponse struct {
	RequestID string          `json:"RequestId,omitempty"`
	User      AssumedRoleUser `json:"AssumedRoleUser"`
	Cred      Credentials     `json:"Credentials"`
}

// AssumeRoleParam is the param for AssumeRole.
type AssumeRoleParam struct {
	RoleArn         string
	RoleSessionName string
	Policy          string
	UID             string
}

// api provides the common methods for a STS API.
// It implements the aliyun.API interface.
type api struct {
	v url.Values
}

func (api) Version() string {
	return Ver
}

func (api) Nonce() string {
	// Nonce can be the same when retry
	return aliyun.Nonce(32)
}

func (api) Method() string {
	return http.MethodGet
}

func (a *api) Param() url.Values {
	return a.v
}

// GetRoleArn composes the roleArn parameter;
// it must be of the form "acs:ram::$accountID:role/$roleName".
func GetRoleArn(uid, roleName string) string {
	return "acs:ram::" + uid + ":role/" + roleName
}

// AssumeRoleAPI forms the API for AssumeRole.
// doc https://help.aliyun.com/document_detail/28763.html
func AssumeRoleAPI(r *AssumeRoleParam, dur int64) aliyun.API {
	a := &api{v: url.Values{}}

	// api-specific mandotory params
	a.v.Add("Action", "AssumeRole")
	a.v.Add("RoleArn", r.RoleArn)
	a.v.Add("RoleSessionName", r.RoleSessionName)

	// optional
	if r.Policy != "" {
		a.v.Add("Policy", r.Policy)
	}
	if dur > 900 && dur < 3600 {
		// default (and max) 3600
		a.v.Add("DurationSeconds", strconv.FormatInt(dur, 10))
	}

	return a
}

// Getter gets the credentials.
type Getter interface {
	// Get gets the credentials using the param.
	// The returned credentials should not expire before
	// specified duration (in seconds).
	Get(*AssumeRoleParam, int64) (Credentials, error)
}

type getter struct {
	s    aliyun.Signer
	host string
	cl   *http.Client
}

func (g *getter) Get(r *AssumeRoleParam, dur int64) (cred Credentials, err error) {
	api := AssumeRoleAPI(r, dur)
	var resp AssumeRoleResponse
	if err = aliyun.Get(g.cl, g.s, api, g.host, &resp); err != nil {
		return
	}
	cred = resp.Cred
	return
}

// New returns a Getter for requesting credentials from host
// with the provided Signer. The underlying http.Client is
// set to have a 5s timeout.
func New(s aliyun.Signer, host string) Getter {
	return &getter{
		s:    s,
		host: host,
		cl: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
