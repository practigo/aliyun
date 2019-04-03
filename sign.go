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

// Signer signs an API.
type Signer interface {
	// Sign signs the API according to the Aliyun specification:
	// https://help.aliyun.com/document_detail/50284.html &
	// https://help.aliyun.com/document_detail/50286.html.
	//
	// Version	String	是	API 版本号。
	// 		格式：YYYY-MM-DD。
	// 		本版本对应为 2016-11-01。
	// AccessKeyId	String	是	阿里云颁发给用户的访问服务所用的密钥 ID。
	// Signature	String	是	签名结果串。关于签名的计算方法，参见 签名机制。
	// SignatureMethod	String	是	签名方式，目前支持 HMAC-SHA1。
	// Timestamp	String	是	请求的时间戳。
	// 		日期格式按照 ISO8601 标准表示，并需要使用 UTC 时间。
	// 		格式：YYYY-MM-DDThh:mm:ssZ。
	// 		例如：2014-05-26T12:00:00Z（为北京时间 2014 年 5 月 26 日 20 点 0 分 0 秒）。
	// SignatureVersion	String	是	签名算法版本。目前版本是 1.0。
	// SignatureNonce	String	是	唯一随机数，用于防止网络重放攻击。用户在不同请求间要使用不同的随机数值。
	// ResourceOwnerAccount	String	否	本次 API 请求访问到的资源拥有者账户，即登录用户名。 此参数的使用方法，参见 RAM资源授权。（只能在 RAM 中可对 live 资源进行授权的 Action 中才能使用此参数，否则访问会被拒绝。）
	// Format	String	否	返回值的类型，支持 JSON 与 XML。默认值：XML
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

// NewSigner returns a Signer to sign the API.
// It is fixed to use the JSON format.
func NewSigner(id, secret string) Signer {
	return &signer{
		id:     id,
		secret: secret,
		ver:    "1.0",
		method: "HMAC-SHA1",
		format: JSONFormat, // force JSON, opinionated
	}
}

// RandString returns a random string with given length.
func RandString(n int) string {
	b := make([]byte, n/2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
