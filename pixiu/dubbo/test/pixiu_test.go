package test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestPost(t *testing.T) {

	url := "http://localhost:8883/api/v1/test-dubbo/UserProvider/com.dubbogo.pixiu.UserService?" +
		"group=dubbo-test&version=1.0.0&method=GetUserByName"
	data := "{\"types\":\"string\",\"values\":\"tc\" }"

	//ctx := context.Background()
	//tpsNum,_ :=strconv.Atoi(os.Getenv("tps"))
	//parallel,_ := strconv.Atoi(os.Getenv("parallel"))

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(s), "0001"))
}
