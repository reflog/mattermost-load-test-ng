package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/mattermost/mattermost-load-test-ng/config"
)

func TestAPI(t *testing.T) {
	// create http.Handler
	handler := SetupAPIRouter()

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// is it working?
	e.GET("/loadtest/status/123").
		Expect().
		Status(http.StatusNotFound)

	sampleConfigBytes, _ := ioutil.ReadFile("../config/config.default.json")
	var sampleConfig config.LoadTestConfig
	_ = json.Unmarshal(sampleConfigBytes, &sampleConfig)
	sampleConfig.ConnectionConfiguration.ServerURL = "http://fakesitetotallydoesntexist.com"
	sampleConfig.UsersConfiguration.MaxActiveUsers = 100
	obj := e.POST("/loadtest/create").WithJSON(sampleConfig).
		Expect().
		Status(http.StatusOK).JSON().Object()
	ltId := obj.Value("loadTestId").String().Raw()

	e.POST("/loadtest/run/" + ltId).Expect().Status(http.StatusOK)
	e.PUT("/loadtest/user/"+ltId).WithQuery("amount", 10).Expect().Status(http.StatusOK)
	e.DELETE("/loadtest/user/"+ltId).WithQuery("amount", 3).Expect().Status(http.StatusOK)
	e.POST("/loadtest/stop/" + ltId).Expect().Status(http.StatusOK)
	e.POST("/loadtest/destroy/" + ltId).Expect().Status(http.StatusOK)
}
