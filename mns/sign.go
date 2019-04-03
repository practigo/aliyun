package mns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

// Signer signs the MNS API.
type Signer interface {
	// Sign returns a Authorization string for the specified API
	// HTTP method, request resource and headers.
	Sign(method, resource string, headers map[string]string) string
}

type signer struct {
	key    string
	secret string
}

// NewSigner returns a Signer to sign the MNS APIs.
func NewSigner(key, secret string) Signer {
	return &signer{
		key:    key,
		secret: secret,
	}
}

// Sign signs the API, as described in
// https://help.aliyun.com/document_detail/27487.html.
//
// Authorization: MNS AccessKeyId:Signature
// Signature = base64(hmac-sha1(HTTP_METHOD + "\n"
// + CONTENT-MD5 + "\n"
// + CONTENT-TYPE + "\n"
// + DATE + "\n"
// + CanonicalizedMNSHeaders
// + CanonicalizedResource))
func (a *signer) Sign(method, resource string, headers map[string]string) string {
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
	return fmt.Sprintf("MNS %s:%s", a.key, signature)
}
