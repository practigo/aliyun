package aliyun

import (
	"net/http"
	"net/url"
	"time"
)

// Aliyun constants
const (
	XMLFormat  = "XML"                  // default by Aliyun official
	JSONFormat = "JSON"                 // force by this package
	TimeFormat = "2006-01-02T15:04:05Z" // official UTC time format
)

// FormatT converts a Time to Aliyun format time string,
// e.g., 2006-01-02T15:04:05Z.
func FormatT(t time.Time) string {
	return t.UTC().Format(TimeFormat)
}

// An API abstracts the unique parts of an Aliyun API
// for a single request. The Signer then signs the API
// to protect against different attacks.
type API interface {
	// Param returns the API specific parameters.
	// These will combine with the common param to get a signature.
	Param() url.Values
	// Version returns the API specific version,
	// usually in the form of a date string.
	Version() string
	// Nonce returns a random string to migate repeat-attack.
	// It gives the freedom for each API to specify it's
	// own randomness.
	Nonce() string
	// Method returns the HTTP method the API used to
	// make the request.
	Method() string
}

// Base implements parts for the API with a 32 bytes
// string nonce and HTTP GET method.
type Base struct{}

// Nonce returns a random 32 bytes string.
func (Base) Nonce() string {
	return RandString(32)
}

// Method returns GET (the most common method).
func (Base) Method() string {
	return http.MethodGet
}

// An OSS defines a unique OSS resource.
// OSS is the fundamental service used by
// most other services.
type OSS struct {
	Bucket   string `json:"OssBucket"`
	Endpoint string `json:"OssEndpoint"`
	Object   string `json:"OssObject"`
}

// FillOSS helps filling the OSS to the param.
func FillOSS(v url.Values, o OSS) {
	v.Set("OssEndpoint", o.Endpoint)
	v.Set("OssBucket", o.Bucket)
	v.Set("OssObject", o.Object)
}
