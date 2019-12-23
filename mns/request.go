package mns

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/practigo/aliyun"
)

// IsNoMessage checks if the err is MessageNotExist,
// which is usually not an error.
func IsNoMessage(err error) bool {
	ce, ok := err.(*aliyun.CanonicalizedError)
	return ok && ce.Code == "MessageNotExist"
}

// service constants
const (
	Version     = "2015-06-06"
	ContentType = "text/xml"
	// headers
	HeaderDate        = "Date"
	HeaderAuth        = "Authorization"
	HeaderContentType = "Content-Type"
	HeaderMD5         = "Content-MD5"
	HeaderVersion     = "x-mns-version"
)

// CommonHeader returns a map for mandotory headers.
// For all common headers, see
// https://help.aliyun.com/document_detail/27485.html.
func CommonHeader() map[string]string {
	headers := make(map[string]string)
	headers[HeaderContentType] = ContentType
	headers[HeaderDate] = time.Now().UTC().Format(http.TimeFormat)
	headers[HeaderVersion] = Version
	return headers
}

// API wraps the HTTP method, request resource (path & query),
// request body and API-wise headers if any. Both Body and
// Headers can be nil if not necessary.
type API struct {
	// HTTP method
	Method string
	// path with query
	Resource string
	// request body
	Body interface{}
	// extra request headers
	Headers map[string]string
}

// Req sends a the API to host and marshal response body
// to resp using the provided http.Client and Signer.
// Set resp to nil if no response body expected.
// NOTE that host is in the form of
// "http(s)://$AccountId.$Region.aliyuncs.com".
func Req(cl *http.Client, s Signer, host string, a *API, resp interface{}) (err error) {
	// body
	content := []byte{}
	if a.Body != nil {
		content, err = xml.Marshal(a.Body)
		if err != nil {
			return
		}
	}

	hs := CommonHeader()
	for k, v := range a.Headers {
		hs[k] = v
	}
	// TODO: optional md5 header
	// md5Str := fmt.Sprintf("%x", md5.Sum(content))
	// TODO: optional host header for HTTP 1.1?
	hs[HeaderAuth] = s.Sign(a.Method, a.Resource, hs)

	// do request
	uri := fmt.Sprintf("%s%s", host, a.Resource)
	req, err := http.NewRequest(a.Method, uri, bytes.NewBuffer(content))
	if err != nil {
		return
	}

	// assign headers
	for k, v := range hs {
		req.Header.Set(k, v)
	}

	rawResp, err := cl.Do(req)
	if err != nil {
		return
	}
	defer rawResp.Body.Close()

	return aliyun.HandleResp(rawResp, xml.Unmarshal, resp)
}
