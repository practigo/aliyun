package mns_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/practigo/aliyun/mns"
)

var (
	// need real AK to test
	testID       = ""
	testSecret   = ""
	testEndpoint = ""
	testQueue    = ""
)

func TestMessager(t *testing.T) {
	s := mns.NewSigner(testID, testSecret)
	messager := mns.NewMessager(s, testEndpoint)

	// send
	resp, err := messager.Send(testQueue, &mns.SendMessageRequest{
		// MessageBody: []byte("hello world"),
		MessageBody: mns.Encode2Base64([]byte("hello world base64")),
	})
	if err != nil {
		t.Error(err)
	} else {
		bs, _ := json.Marshal(resp)
		t.Log(string(bs))
	}

	time.Sleep(5 * time.Second)

	// receive
	resp2, err := messager.Receive(testQueue, mns.MaxWaitSeconds)
	if err != nil {
		if mns.IsNoMessage(err) {
			t.Log("no message")
		} else {
			t.Error(err)
		}
	} else {
		bs, _ := json.Marshal(resp2)
		t.Log(string(bs))
		body, _ := base64.StdEncoding.DecodeString(string(resp2.MessageBody))
		t.Log(string(body))

		// delete
		err = messager.Delete(testQueue, resp2.ReceiptHandle)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(resp2.ReceiptHandle, "deleted")
		}
	}
}
