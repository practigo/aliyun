package aliyun

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// HandleResp reads the raw HTTP response and checks if
// any error returned. Errors other than io/json error
// should be a CanonicalizedError. The body is unmarshaled
// to the provided interface resp if the request succeeds.
// It is the caller's responsiblity to close the body.
func HandleResp(raw *http.Response, resp interface{}) (err error) {
	bs, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		return
	}

	// try error first
	var ce CanonicalizedError
	err = json.Unmarshal(bs, &ce)
	if err != nil {
		return
	}

	ce.Status = raw.StatusCode

	if ce.Code != "" { // if there's an error, there should be a code
		return &ce
	}

	err = json.Unmarshal(bs, resp)
	return
}

// Get makes a HTTP GET request to the host
// for the provided API, and marshal the response
// into the value pointed by resp.
func Get(cl *http.Client, s Signer, a API, host string, resp interface{}) error {
	url := host + "?" + s.Sign(a)
	r, err := cl.Get(url)
	if r != nil && r.Body != nil {
		defer r.Body.Close()
	}
	if err != nil {
		return err
	}
	return HandleResp(r, resp)
}
