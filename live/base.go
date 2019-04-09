package live

import (
	"net/url"

	"github.com/practigo/aliyun"
)

// API constants
const (
	Host = "https://live.aliyuncs.com"
	Ver  = "2016-11-01"
)

// An api provides the common parts for a Live API.
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

// A StreamURI defines a RTMP stream,
// e.g., rtmp://{Domain}/{App}/{Stream}.
type StreamURI struct {
	Domain string `json:"DomainName"`
	App    string `json:"AppName"`
	Stream string `json:"StreamName"`
}

func fillURI(v url.Values, uri StreamURI) {
	v.Set("DomainName", uri.Domain)
	v.Set("AppName", uri.App)
	v.Set("StreamName", uri.Stream)
}
