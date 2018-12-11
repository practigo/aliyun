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
	Host = "https://sts.aliyuncs.com/"
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

// api provides the common methods for a STS API.
// It implements the aliyun.API interface.
type api struct {
	v url.Values
}

func (api) Version() string {
	return Ver
}

func (api) Nonce() string {
	// TODO: Nonce can be the same when retry
	return aliyun.Nonce(32)
}

func (api) Method() string {
	return http.MethodGet
}

func (a *api) Param() url.Values {
	return a.v
}

func getRoleArn(uid string, role string) string {
	return "acs:ram::" + uid + ":role/" + role
}

// AssumeRole gets a temporary role.
// doc https://help.aliyun.com/document_detail/28763.html
func AssumeRole(s aliyun.Signer, uid, role, sessionName, policy string, dur uint64) (resp AssumeRoleResponse, err error) {
	a := &api{v: url.Values{}}
	// api-specific
	a.v.Add("Action", "AssumeRole")
	a.v.Add("RoleArn", getRoleArn(uid, role))
	a.v.Add("RoleSessionName", sessionName)
	a.v.Add("DurationSeconds", strconv.FormatUint(dur, 10))
	if policy != "" {
		a.v.Add("Policy", policy)
	}
	// request
	url := Host + "?" + s.Sign(a)
	err = aliyun.Get(url, &resp)
	return
}
