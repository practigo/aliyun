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

	s := aliyun.NewAccessKey(testEnvs["STS_KEY_ID"], testEnvs["STS_KEY_SECRET"])
	role := sts.GetRoleArn(testEnvs["STS_UID"], testEnvs["STS_ROLE"])
	param := sts.AssumeRoleParam{
		RoleArn:         role,
		RoleSessionName: testEnvs["STS_SECCSION"],
	}
	g := sts.New(s, sts.Host)
	cred, err := g.Get(&param, 0)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", cred)
	}

	param.RoleArn = "dummy" // make it an error
	_, err = g.Get(&param, 0)
	if err == nil {
		t.Error("should have error with dummy role")
	} else {
		t.Log(err)
	}
}

func TestCache(t *testing.T) {
	if !completed {
		t.Skip("Must set env vars", requiredVars)
	}

	s := aliyun.NewAccessKey(testEnvs["STS_KEY_ID"], testEnvs["STS_KEY_SECRET"])
	role := sts.GetRoleArn(testEnvs["STS_UID"], testEnvs["STS_ROLE"])
	param := sts.AssumeRoleParam{
		RoleArn:         role,
		RoleSessionName: testEnvs["STS_SECCSION"],
	}
	cache := sts.Wrap(sts.New(s, sts.Host), sts.DefaultKey)
	cred, err := cache.Get(&param, 900)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", cred)
	}
}
