# Go SDK for [Aliyun](https://www.aliyun.com)

[![GoDoc](https://godoc.org/github.com/practigo/aliyun?status.svg)](https://godoc.org/github.com/practigo/aliyun)
[![Go Report Card](https://goreportcard.com/badge/github.com/practigo/aliyun)](https://goreportcard.com/report/github.com/practigo/aliyun)

This is not a typical SDK as it won't try to provide all APIs. 

For the complete APIs see the official https://github.com/aliyun/alibaba-cloud-sdk-go.

## Products \& APIs

### [STS](https://help.aliyun.com/document_detail/28756.html)

- AssumeRole
- Cache

### [MTS](https://help.aliyun.com/document_detail/66804.html)

- Transcode-Job

### Live

- Record create & decribe

Other repo: https://github.com/BPing/aliyun-live-go-sdk

### OSS

https://github.com/aliyun/aliyun-oss-go-sdk is good enough.

### MNS

- Different Restful APIs and authorization
- Queue & Messages

## Usage

The typical usage is:

```go
s := aliyun.NewSigner("id", "secret")
api := YourImplementedAPI{}
url := "host" + ? + s.Sign(api)
rawResponse, err := http.Get(url) // http level error if any
var resp APIResponseType
err = HandleResp(rawResponse, &resp) // return a CanonicalizedError if any
```

## License

MIT
