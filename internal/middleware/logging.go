package middleware

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LogLine struct {
	Timestamp  string  `json:"ts"`
	Tool       string  `json:"tool"`
	Function   string  `json:"fn"`
	LatencyMS  float64 `json:"latency_ms"`
	IsOK       bool    `json:"ok"`
	StatusCode int     `json:"http_status"`
}

func JsonLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		pathParts := strings.Split(strings.Trim(c.Request.URL.Path, "/"), "/")
		var tool, fn string
		if len(pathParts) >= 3 {
			tool = pathParts[1]
			fn = pathParts[2]
		} else {
			tool = "unknown"
			fn = "unknown"
		}

		statusCode := c.Writer.Status()
		isOK := statusCode < 400

		logData := &LogLine{
			Timestamp:  end.UTC().Format(time.RFC3339Nano),
			Tool:       tool,
			Function:   fn,
			LatencyMS:  float64(latency.Microseconds()) / 1000.0,
			IsOK:       isOK,
			StatusCode: statusCode,
		}

		logLine, err := json.Marshal(logData)
		if err != nil {
			log.Printf("Error marshaling JSON log: %v", err)
			return
		}

		log.Println(string(logLine))
	}
}
