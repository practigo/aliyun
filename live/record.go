package live

import (
	"net/http"
	"net/url"
	"time"

	"github.com/practigo/aliyun"
)

// DescribeRecordsAPI returns the API for DescribeLiveStreamRecordIndexFiles.
// https://help.aliyun.com/document_detail/35423.html
func DescribeRecordsAPI(uri StreamURI, start, end time.Time) aliyun.API {
	a := &api{v: url.Values{}}

	// api-specific
	a.v.Add("Action", "DescribeLiveStreamRecordIndexFiles")
	fillURI(a.v, uri)
	a.v.Add("StartTime", aliyun.FormatT(start))
	a.v.Add("EndTime", aliyun.FormatT(end))

	return a
}

// CreateRecordAPI returns the API for CreateLiveStreamRecordIndexFiles.
// https://help.aliyun.com/document_detail/35417.html
func CreateRecordAPI(uri StreamURI, start, end time.Time, oss aliyun.OSS) aliyun.API {
	a := &api{v: url.Values{}}

	// api-specific
	a.v.Add("Action", "CreateLiveStreamRecordIndexFiles")
	fillURI(a.v, uri)
	aliyun.FillOSS(a.v, oss)
	a.v.Add("StartTime", aliyun.FormatT(start))
	a.v.Add("EndTime", aliyun.FormatT(end))

	return a
}

// A RecordInfo represents the a live stream record.
type RecordInfo struct {
	RecordID  string `json:"RecordId"`
	RecordURL string `json:"RecordUrl"`
	// timestamps
	CreateTime time.Time `json:"CreateTime"`
	StartTime  time.Time `json:"StartTime"`
	EndTime    time.Time `json:"EndTime"`
	// video related
	Duration float64 `json:"Duration"`
	Width    int     `json:"Width"`
	Height   int     `json:"Height"`
	// embedded
	StreamURI
	aliyun.OSS
}

// A DescribeRecordsResponse is the response for
// DescribeLiveStreamRecordIndexFiles.
type DescribeRecordsResponse struct {
	List struct {
		Files []RecordInfo `json:"RecordIndexInfo"`
	} `json:"RecordIndexInfoList"`
	RequestID string `json:"RequestId"`
}

// A CreateRecordResponse is the response for
// CreateLiveStreamRecordIndexFiles.
type CreateRecordResponse struct {
	Info      RecordInfo `json:"RecordInfo"`
	RequestID string     `json:"RequestId"`
}

// DescribeRecords uses the signer to send a DescribeRecordsAPI.
func DescribeRecords(s aliyun.Signer, uri StreamURI, start, end time.Time) (resp DescribeRecordsResponse, err error) {
	api := DescribeRecordsAPI(uri, start, end)

	err = aliyun.Get(http.DefaultClient, s, api, Host, &resp)
	return
}

// CreateRecord uses the signer to send a CreateRecordAPI.
func CreateRecord(s aliyun.Signer, uri StreamURI, start, end time.Time, oss aliyun.OSS) (resp CreateRecordResponse, err error) {
	api := CreateRecordAPI(uri, start, end, oss)

	err = aliyun.Get(http.DefaultClient, s, api, Host, &resp)
	return
}
