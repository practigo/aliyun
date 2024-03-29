package mts_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/practigo/aliyun"
	"github.com/practigo/aliyun/mts"
)

var (
	testEnvs     = make(map[string]string)
	requiredVars = []string{"MTS_KEY_ID", "MTS_KEY_SECRET", "MTS_ENDPOINT"}
	completed    = true

	tr mts.Transcoder
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

	s := aliyun.NewAccessKey(testEnvs["MTS_KEY_ID"], testEnvs["MTS_KEY_SECRET"])
	tr = mts.New(s, testEnvs["MTS_ENDPOINT"])

	os.Exit(m.Run())
}

func TestSubmit(t *testing.T) {
	if !completed {
		t.Skip("Must set env vars", requiredVars)
	}

	// CHANGE the request to make it non-error.
	// Dummy example from official doc.
	req := mts.SubmitJobsRequest{
		Input:          `{"Bucket":"example-bucket","Location":"oss-cn-hangzhou","Object":"example.flv"}`,
		Outputs:        `[{"OutputObject":"example-output.flv","TemplateId":"S00000000-000010","WaterMarks":[{"InputFile":{"Bucket":"example-bucket","Location":"oss-cn-hangzhou","Object":"example-logo.png"},"WaterMarkTemplateId":"88c6ca184c0e47098a5b665e2a126797"}],"UserData":"testid-001"}]`,
		OutputBucket:   "example-bucket",
		OutputLocation: "oss-cn-shanghai",
		PipelineID:     "example-pipeline",
	}
	resp, err := tr.Submit(&req)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("%+v", resp)
	}
}

func TestQuery(t *testing.T) {
	if !completed {
		t.Skip("Must set env vars", requiredVars)
	}

	// CHANGE the request to make it non-error.
	resp2, err := tr.Query("example-jobid1", "example-jobid2")
	if err != nil {
		t.Log(err)
	} else {
		bs, _ := json.MarshalIndent(resp2, "", "  ")
		t.Log(string(bs))
	}
}
