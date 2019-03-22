/*
Package aliyun contains a library for development around Alibaba Cloud Services.

It aims to provide primitives for making requests to any services,
and wrappers/tools around APIs for better user-friendly experience.

The API, Signer & CanonicalizedError are three key primitives that
can be used for almost all services.

The typical usage is:

```go
s := aliyun.NewSigner("id", "secret")
api := YourImplementedAPI{}
url := "host" + ? + s.Sign(api)
rawResponse, err := http.Get(url) // http level error if any
var resp APIResponseType
err = HandleResp(rawResponse, &resp) // return a CanonicalizedError if any
```

See sub-packages for API examples.
*/
package aliyun
