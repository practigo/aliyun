package aliyun

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// HandleResp reads the HTTP response and checks if any error returned.
// If any error other than io/json error returned, it should be a
// CanonicalizedError. The body is unmarshaled to the provided interface v
// if the request succeeds.
func HandleResp(r *http.Response, v interface{}) (err error) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// try error first
	var ce CanonicalizedError
	err = json.Unmarshal(bs, &ce)
	if err != nil {
		return
	}

	ce.Status = r.StatusCode

	if ce.Code != "" { // if there's an error, there should be a code
		return &ce
	}

	err = json.Unmarshal(bs, &v)
	return
}

// Get makes a HTTP GET request to the provided url
// and handles the response.
func Get(url string, v interface{}) error {
	r, err := http.Get(url)
	if r != nil {
		defer r.Body.Close()
	}
	if err != nil {
		return err
	}
	return HandleResp(r, &v)
}
