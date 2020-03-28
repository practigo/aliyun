package acm_test

import (
	"os"
	"testing"
	"time"

	"github.com/practigo/aliyun/acm"
)

var (
	testEnvs     = make(map[string]string)
	requiredVars = []string{"ACM_KEY_ID", "ACM_KEY_SECRET", "ACM_GROUP", "ACM_DATA_ID"}
	completed    = true
	// optional
	ak   *acm.AccessKey
	opt  acm.ConfigOption
	host string
)

func TestMain(m *testing.M) {
	for _, k := range requiredVars {
		if v := os.Getenv(k); v != "" {
			testEnvs[k] = v
		} else {
			completed = false
			break
		}
	}
	host = os.Getenv("ACM_HOST")
	if host == "" {
		host = "acm.aliyun.com" // public domain
	}
	if completed {
		ak = &acm.AccessKey{
			AccessKeyID:     testEnvs["ACM_KEY_ID"],
			AccessKeySecret: testEnvs["ACM_KEY_SECRET"],
			SecurityToken:   os.Getenv("ACM_TOKEN"),
		}
		opt = acm.ConfigOption{
			Group:  testEnvs["ACM_GROUP"],
			DataID: testEnvs["ACM_DATA_ID"],
			Tenant: os.Getenv("ACM_TENANT"),
		}
	}
	os.Exit(m.Run())
}

func TestSignature(t *testing.T) {
	abc := acm.AccessKey{
		AccessKeySecret: "abc",
	}
	s := abc.Sign("test-content")
	if s != "9bPdA2uiklekqbIziaQnjEEphf4=" {
		t.Error("wrong sign:", s)
	}
}

func TestGetIPs(t *testing.T) {
	srv := acm.New("acm.aliyun.com")
	t.Log(srv.GetIPs())
}

func TestService(t *testing.T) {
	if !completed {
		t.Skip("Must set env vars", requiredVars)
	}
	srv := acm.New(host)
	resp, err := srv.GetConfig(ak, opt)
	if err != nil {
		t.Fatal("get config:", err)
	}
	t.Log(string(resp))
}

func TestListenConfig(t *testing.T) {
	if !completed {
		t.Skip("Must set env vars", requiredVars)
	}
	srv := acm.New(host)
	resp, err := srv.ListenConfig(ak, opt)
	if err != nil {
		t.Fatal("listen config:", err)
	}
	opts := acm.ParseListenResponse(resp)
	if len(opts) < 1 || opts[0].DataID != opt.DataID {
		t.Fatal("mismatch dataID", string(resp), opts)
	}

	// get
	resp, err = srv.GetConfig(ak, opt)
	if err != nil {
		t.Fatal("get config:", err)
	}

	md5 := acm.MD5(resp)
	opt.MD5 = md5
	st := time.Now()
	resp, err = srv.ListenConfig(ak, opt)
	if err != nil {
		t.Fatal("poll config:", err)
	}
	if len(resp) > 1 {
		t.Error("should be empty resp:", string(resp))
	}
	dur := time.Now().Sub(st)
	if dur < 25*time.Second {
		t.Error("should be 30s long polling:", dur)
	}

	t.Log(md5, dur)
}
