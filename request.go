package aliyun

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// HandleResp reads the raw HTTP response and checks if
// any error returned. Errors other than io/json error
// should be a CanonicalizedError. The body is unmarshaled
// to the provided interface resp using the function f.
// If f is nil, it defaults to json.Unmarshal.
// It is the caller's responsibility to close the body.
func HandleResp(raw *http.Response, f func([]byte, interface{}) error, resp interface{}) (err error) {
	if f == nil {
		f = json.Unmarshal // default
	}

	bs, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		return
	}

	// try error first
	var ce CanonicalizedError
	err = f(bs, &ce)
	if err != nil {
		return
	}

	ce.Status = raw.StatusCode

	// TODO: check status also?
	// if ce.Code != "" || ce.Status < 200 || ce.Status > 299 {
	if ce.Code != "" { // if there's an error, there should be a code
		return &ce
	}

	err = f(bs, resp)
	return
}

// Get makes a HTTP GET request to the host
// for the provided API, and marshal the response
// into the value pointed by resp.
func Get(cl *http.Client, s Signer, a API, host string, resp interface{}) error {
	url := host + "?" + s.Sign(a)
	r, err := cl.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return HandleResp(r, json.Unmarshal, resp)
}
