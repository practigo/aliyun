package aliyun

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

// A Signer signs the APIs.
type Signer interface {
	// Sign signs the API according to the Aliyun specification:
	// https://help.aliyun.com/document_detail/50284.html &
	// https://help.aliyun.com/document_detail/50286.html.
	// It returns the query part of the request including
	// the signature.
	Sign(API) string
}

type signer struct {
	// user specific
	id, secret string
	// internal, so far fixed
	ver    string
	method string
	format string
}

func (s *signer) Sign(a API) string {
	// API specific param
	v := a.Param()

	// public param (overwrite if already exists, which should not)
	v.Set("Version", a.Version())
	v.Set("AccessKeyId", s.id)
	v.Set("SignatureMethod", s.method)
	v.Set("Timestamp", FormatT(time.Now()))
	v.Set("SignatureVersion", s.ver)
	v.Set("SignatureNonce", a.Nonce())
	v.Set("Format", s.format)

	// this also sort the params
	query := v.Encode()

	// StringToSign=
	// HTTPMethod + “&” +
	// percentEncode(“/”) + ”&” +
	// percentEncode(CanonicalizedQueryString)
	toSign := a.Method() + "&%2F&" + url.QueryEscape(query)

	// generate signature
	h := hmac.New(sha1.New, []byte(s.secret+"&"))
	h.Write([]byte(toSign)) // sha1 Write() returns no error
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// final query
	return query + "&Signature=" + url.QueryEscape(signature)
	// or
	// v.Set("Signature", signature)
	// return v.Encode()
}

// NewSigner returns a Signer to sign the APIs
// using the JSON format.
func NewSigner(id, secret string) Signer {
	return &signer{
		id:     id,
		secret: secret,
		ver:    "1.0",
		method: "HMAC-SHA1",
		format: JSONFormat, // force JSON, opinionated
	}
}

// RandString returns a random string with the given
// length n.
func RandString(n int) string {
	b := make([]byte, n/2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
