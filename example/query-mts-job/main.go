package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/practigo/aliyun"
	"github.com/practigo/aliyun/mts"
)

func main() {
	var (
		testEnvs     = make(map[string]string)
		requiredVars = []string{"MTS_KEY_ID", "MTS_KEY_SECRET", "MTS_ENDPOINT"}
	)
	for _, k := range requiredVars {
		if v := os.Getenv(k); v != "" {
			testEnvs[k] = v
		} else {
			panic(fmt.Sprintf("require env setting: %v", requiredVars))
		}
	}

	s := aliyun.NewAccessKey(testEnvs["MTS_KEY_ID"], testEnvs["MTS_KEY_SECRET"])
	tr := mts.New(s, testEnvs["MTS_ENDPOINT"])

	resp, err := tr.Query(os.Args[1], os.Args[2:]...)
	if err != nil {
		fmt.Println(err)
	} else {
		bs, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(bs))
	}
}
