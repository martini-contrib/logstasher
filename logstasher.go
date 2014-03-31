// Package logstasher is a Martini middleware that prints logstash-compatiable
// JSON to a given io.Writer for each HTTP request.
package logstasher

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-martini/martini"
)

type logstashEvent struct {
	Timestamp string              `json:"@timestamp"`
	Version   int                 `json:"@version"`
	Method    string              `json:"method"`
	Path      string              `json:"path"`
	Status    int                 `json:"status"`
	Size      int                 `json:"size"`
	Duration  float64             `json:"duration"`
	Params    map[string][]string `json:"params,omitempty"`
}

// Logger returns a middleware handler prints the request in a Logstash-JSON compatiable format
func Logger(writer io.Writer) martini.Handler {
	out := log.New(writer, "", 0)
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {
		start := time.Now()

		rw := res.(martini.ResponseWriter)
		c.Next()

		event := logstashEvent{
			time.Now().Format(time.RFC3339),
			1,
			req.Method,
			req.URL.Path,
			rw.Status(),
			rw.Size(),
			time.Since(start).Seconds() * 1000.0,
			map[string][]string(req.Form),
		}

		output, err := json.Marshal(event)
		if err != nil {
			// Should this be fatal?
			log.Printf("Unable to JSON-ify our event (%#v): %v", event, err)
			return
		}
		out.Println(string(output))
	}
}
