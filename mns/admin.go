package mns

import (
	"encoding/xml"
	"net/http"
)

// QueueAttributes ...
type QueueAttributes struct {
	XMLName                xml.Name `xml:"Queue" json:"-"`
	QueueName              string   `xml:"QueueName"`
	CreateTime             int64    `xml:"CreateTime"`
	LastModifyTime         int64    `xml:"LastModifyTime"`
	VisibilityTimeout      int64    `xml:"VisibilityTimeout"`
	MaximumMessageSize     int64    `xml:"MaximumMessageSize"`
	MessageRetentionPeriod int64    `xml:"MessageRetentionPeriod"`
	DelaySeconds           int64    `xml:"DelaySeconds"`
	PollingWaitSeconds     int64    `xml:"PollingWaitSeconds"`
	InactiveMessages       int64    `xml:"InactiveMessages"`
	ActiveMessages         int64    `xml:"ActiveMessages"`
	DelayMessages          int64    `xml:"DelayMessages"`
	LoggingEnabled         bool     `xml:"LoggingEnabled"`
}

// GetQueueAttributes returns API for
// https://help.aliyun.com/document_detail/35131.html
func GetQueueAttributes(name string) *API {
	return &API{
		Method:   http.MethodGet,
		Resource: "/queues/" + name,
	}
}
