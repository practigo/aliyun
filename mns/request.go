package mns

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/practigo/aliyun"
)

func readAndUnmarshalBody(body io.Reader, v interface{}) error {
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "read resp body")
	}
	if err = xml.Unmarshal(bs, v); err != nil {
		return errors.Wrap(err, "unmarshal resp xml")
	}
	return nil
}

// HandleResp reads the raw HTTP response and checks if
// any error returned. Errors other than io/xml error
// should be a CanonicalizedError. The body is unmarshaled
// to the provided interface resp if the request succeeds
// and resp is not nil.
// It is the caller's responsibility to close the body.
func HandleResp(raw *http.Response, resp interface{}) (err error) {
	status := raw.StatusCode

	if status == http.StatusCreated ||
		status == http.StatusOK ||
		status == http.StatusNoContent {
		// succeed
		if resp != nil {
			// body expected
			if err = readAndUnmarshalBody(raw.Body, resp); err != nil {
				return errors.Wrap(err, "handle successful resp")
			}
		}
		return nil
	}

	// error case
	var ce aliyun.CanonicalizedError
	if err = readAndUnmarshalBody(raw.Body, &ce); err != nil {
		return errors.Wrap(err, "handle error resp")
	}
	ce.Status = status
	return &ce
}

// IsNoMessage checks if the err is MessageNotExist,
// which is usually not an error.
func IsNoMessage(err error) bool {
	if ce, ok := err.(*aliyun.CanonicalizedError); ok {
		return ce.Code == "MessageNotExist"
	}
	return false
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
			err = errors.Wrap(err, "marshal XML")
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

	return HandleResp(rawResp, resp)
}
