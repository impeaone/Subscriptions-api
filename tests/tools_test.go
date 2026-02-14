package tests

import (
	"agrigation_api/pkg/tools"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestEnvTools(t *testing.T) {
	env1 := os.Getenv("ENV1")
	env1Test := tools.GetEnv("ENV1", "")
	if env1 != env1Test {
		t.Error(env1, env1Test)
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	env2 := os.Getenv("ENV2")
	env2Test := tools.GetEnvAsBool("ENV2", false)
	if env2 == "" && env2Test != false {
		t.Error(env2, env2Test)
	}
	if env2 != "" && env2Test != true {
		t.Error(env2, env2Test)
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	env3 := os.Getenv("345")
	env3Test := tools.GetEnvAsInt("345", 0)
	if env3 != strconv.Itoa(env3Test) && env3 != "" {
		t.Error(env1, env1Test)
	}
	if env3 == "" && env3Test != 0 {
		t.Error(env1, env1Test)
	}
}

type TestResponseWriter struct{}

func (w TestResponseWriter) Header() http.Header {
	return http.Header{}
}
func (w TestResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}
func (w TestResponseWriter) WriteHeader(statusCode int) {}

func TestWriteTools(t *testing.T) {
	w := TestResponseWriter{}
	data := "testdata"
	tools.WriteJSON(w, 200, data)
	tools.WriteError(w, 500, data)
}

func TestParseTools(t *testing.T) {
	testTime := "07-2025"
	testParseTime, _ := time.Parse("01-2006", testTime)
	testedTime, _ := tools.ParseMonthYear(testTime)
	if testParseTime != testedTime {
		t.Error(testTime)
	}
}
