package sts_test

import (
	"os"
	"testing"

	"github.com/practigo/aliyun"
	"github.com/practigo/aliyun/sts"
)

var (
	testEnvs     = make(map[string]string)
	requiredVars = []string{"STS_KEY_ID", "STS_KEY_SECRET", "STS_UID", "STS_ROLE", "STS_SECCSION"}
	completed    = true
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
	os.Exit(m.Run())
}

func TestAssumeRole(t *testing.T) {
	if !completed {
		t.Skip("Must set env vars", requiredVars)
	}
	s := aliyun.NewSigner(testEnvs["STS_KEY_ID"], testEnvs["STS_KEY_SECRET"])
	resp, err := sts.AssumeRole(s, testEnvs["STS_UID"], testEnvs["STS_ROLE"], testEnvs["STS_SECCSION"], "", 3600)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", resp)
	}
}
