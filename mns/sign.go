package mns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

// A Signer signs the MNS API.
type Signer interface {
	// Sign returns a Authorization string for the specified API
	// HTTP method, request resource and headers, as described in
	// https://help.aliyun.com/document_detail/27487.html.
	//
	// The return string is in the form of:
	// Authorization: MNS AccessKeyId:Signature
	//
	// Signature = base64(hmac-sha1(HTTP_METHOD + "\n"
	// + CONTENT-MD5 + "\n"
	// + CONTENT-TYPE + "\n"
	// + DATE + "\n"
	// + CanonicalizedMNSHeaders
	// + CanonicalizedResource))
	Sign(method, resource string, headers map[string]string) string
}

// AK is the access key.
type AK struct {
	id     string
	secret string
}

// NewAK returns an AK to sign the MNS APIs.
func NewAK(id, secret string) *AK {
	return &AK{
		id:     id,
		secret: secret,
	}
}

// Sign returns the Authorization string.
func (a *AK) Sign(method, resource string, headers map[string]string) string {
	// CanonicalizedMNSHeaders
	mnsHeaders := []string{}
	for k, v := range headers {
		if strings.HasPrefix(k, "x-mns-") {
			mnsHeaders = append(mnsHeaders, k+":"+strings.TrimSpace(v))
		}
	}
	sort.Sort(sort.StringSlice(mnsHeaders))

	toSign := method + "\n" +
		headers[HeaderMD5] + "\n" +
		headers[HeaderContentType] + "\n" +
		headers[HeaderDate] + "\n" +
		strings.Join(mnsHeaders, "\n") + "\n" +
		resource

	h := hmac.New(sha1.New, []byte(a.secret))
	h.Write([]byte(toSign)) // no error here

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("MNS %s:%s", a.id, signature)
}
