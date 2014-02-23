package logstash_logger

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/codegangsta/martini"
)

type logstashEvent struct {
	Timestamp string  `json:"@timestamp"`
	Version   int     `json:"@version"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	Status    int     `json:"status"`
	Size      int     `json:"size"`
	Duration  float64 `json:"duration"`
}

// Logger returns a middleware handler prints the request in a Logstash-JSON compatiable format
func Logger(out *log.Logger) martini.Handler {
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
		}

		output, err := json.Marshal(event)
		if err == nil {
			out.Println(string(output))
		} else {
			// Should this be fatal?
			log.Printf("Unable to JSON-ify our event (%#v): %v", event, err)
		}
	}
}
