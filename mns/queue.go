package mns

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

// SendMessageRequest is the request body for API
// SendMessage. Note that the MessageBody is raw
// []byte here. Whether using base64 encoding as
// other SDKs do is up to the users. See
// Encode2Base64() for reference.
type SendMessageRequest struct {
	XMLName      xml.Name `xml:"Message"`
	MessageBody  []byte   `xml:"MessageBody"`
	DelaySeconds int32    `xml:"DelaySeconds,omitempty"`
	Priority     int32    `xml:"Priority,omitempty"`
}

// SendMessageResponse is the response body for API
// SendMessage.
type SendMessageResponse struct {
	XMLName        xml.Name `xml:"Message" json:"-"`
	MessageID      string   `xml:"MessageId" json:"message_id"`
	MessageBodyMD5 string   `xml:"MessageBodyMD5" json:"message_body_md5"`
	ReceiptHandle  string   `xml:"ReceiptHandle,omitempty" json:"receipt_handle,omitempty"`
}

// ReceiveMessageResponse is the response body for
// API ReceiveMessage.
type ReceiveMessageResponse struct {
	XMLName          xml.Name `xml:"Message" json:"-"`
	MessageID        string   `xml:"MessageId" json:"message_id"`
	ReceiptHandle    string   `xml:"ReceiptHandle" json:"receipt_handle"`
	MessageBodyMD5   string   `xml:"MessageBodyMD5" json:"message_body_md5"`
	MessageBody      []byte   `xml:"MessageBody" json:"message_body"`
	EnqueueTime      int64    `xml:"EnqueueTime" json:"enqueue_time"`
	NextVisibleTime  int64    `xml:"NextVisibleTime" json:"next_visible_time"`
	FirstDequeueTime int64    `xml:"FirstDequeueTime" json:"first_dequeue_time"`
	DequeueCount     int32    `xml:"DequeueCount" json:"dequeue_count"`
	Priority         int32    `xml:"Priority" json:"priority"`
}

// queue constants
const (
	// SendMessage
	DefaultDelay    = 0
	MaxDelay        = 604800
	HighestPriority = 1
	LowestPriority  = 16
	DefaultPriority = 8
	MaxBodyLength   = 64 * 1024

	// BatchSendMessage
	MaxMsgInBatch = 16

	// ReceiveMessage
	MaxWaitSeconds = 30

	// visibility
	MinVisTimeout = 1
	MaxVisTimeout = 43200

	// paths
	queueMsgPath = "/queues/%s/messages"

	// URL queries
	receiveWaitParam = "waitseconds"
	batchNumParam    = "numOfMessages"
	visibilityParam  = "visibilityTimeout"
	peekonlyParam    = "peekonly=true"
)

// Encode2Base64 encodes the src bytes to base64 bytes.
func Encode2Base64(src []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)
	return buf
}

// Messager provides operations on the messages on
// queues of a certain host.
type Messager struct {
	s      Signer
	host   string
	cl     *http.Client
	poller *http.Client
}

// Send sends a message to the queue. See
// https://help.aliyun.com/document_detail/35134.html.
func (m *Messager) Send(queue string, msg *SendMessageRequest) (resp SendMessageResponse, err error) {
	a := &API{
		Method:   http.MethodPost,
		Resource: fmt.Sprintf(queueMsgPath, queue),
		Body:     msg,
	}
	err = Req(m.cl, m.s, m.host, a, &resp)
	return
}

// Receive receives a message from the queue. See
// https://help.aliyun.com/document_detail/35136.html.
func (m *Messager) Receive(queue string, wait int) (resp ReceiveMessageResponse, err error) {
	a := &API{
		Method:   http.MethodGet,
		Resource: fmt.Sprintf(queueMsgPath, queue),
	}
	if wait > 0 {
		a.Resource += fmt.Sprintf("?%s=%d", receiveWaitParam, wait)
	}
	err = Req(m.poller, m.s, m.host, a, &resp)
	return
}

// Delete deletes a message from the queue. See
// https://help.aliyun.com/document_detail/35138.html.
func (m *Messager) Delete(queue string, receipt string) error {
	a := &API{
		Method:   http.MethodDelete,
		Resource: fmt.Sprintf(queueMsgPath+"?ReceiptHandle=%s", queue, receipt),
	}
	return Req(m.cl, m.s, m.host, a, nil)
}

// NewMessager returns a *Messager with the underlying
// http.Clients set.
func NewMessager(s Signer, host string) *Messager {
	return &Messager{
		s:    s,
		host: host,
		cl: &http.Client{
			Timeout: 5 * time.Second,
		},
		poller: &http.Client{
			Timeout: 35 * time.Second,
		},
	}
}
