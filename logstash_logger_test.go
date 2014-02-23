package logstash_logger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"

	"github.com/codegangsta/martini"
)

// Copied and adapted from martini's logger_test.go
func Test_Logger(t *testing.T) {
	buff := bytes.NewBufferString("")
	logger := log.New(buff, "", 0)
	recorder := httptest.NewRecorder()

	m := martini.New()
	m.Use(Logger(logger))
	m.Use(func(res http.ResponseWriter) {
		res.WriteHeader(http.StatusNotFound)
	})

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)

	var event logstashEvent
	err = json.Unmarshal(buff.Bytes(), &event)
	if err != nil {
		t.Error(err)
	}
	expect(t, event.Version, 1)
	expect(t, event.Method, "GET")
	expect(t, event.Path, "/foobar")
	expect(t, event.Status, http.StatusNotFound)
	refute(t, event.Duration, 0)
	refute(t, event.Timestamp, "")
	expect(t, event.Size, 0)
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Expected %v (type %v) - Got %v (type %v)", line, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Did not expect %v (type %v) - Got %v (type %v)", line, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
