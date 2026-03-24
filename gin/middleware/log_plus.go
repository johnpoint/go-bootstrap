package apimiddleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/johnpoint/go-bootstrap/v2/utils"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// sensitiveHeaders lists header names that should be redacted in logs.
var sensitiveHeaders = []string{
	"authorization",
	"cookie",
	"set-cookie",
	"x-api-key",
	"x-auth-token",
	"x-access-token",
}

type reqLog struct {
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
	URL    string      `json:"URL"`
	Resp   string      `json:"resp"`
	ReqID  string      `json:"req-id"`
	In     time.Time   `json:"in"`
	Out    time.Time   `json:"out"`
	Method string      `json:"method"`
}

// LogPlusConfig holds configuration options for LogPlusMiddlewareWithConfig.
type LogPlusConfig struct {
	// LogBody controls whether request and response bodies are logged.
	// WARNING: Enabling this may expose sensitive data. Default: false.
	LogBody bool
	// LogHeaders controls whether request headers are logged.
	// Sensitive headers (Authorization, Cookie, etc.) are always redacted.
	LogHeaders bool
}

// CustomResponseWriter is a wrapper around gin.ResponseWriter that captures response body data.
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write writes the given byte slice to both the buffer and the underlying response writer.
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString writes the given string to both the buffer and the underlying response writer.
func (w *CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// redactHeaders returns a copy of headers with sensitive values replaced by "[REDACTED]".
func redactHeaders(headers http.Header) http.Header {
	redacted := headers.Clone()
	for key := range redacted {
		lower := strings.ToLower(key)
		for _, sensitive := range sensitiveHeaders {
			if lower == sensitive {
				redacted.Set(key, "[REDACTED]")
				break
			}
		}
		if strings.Contains(lower, "token") || strings.Contains(lower, "secret") || strings.Contains(lower, "key") {
			redacted.Set(key, "[REDACTED]")
		}
	}
	return redacted
}

// LogPlusMiddleware returns a Gin middleware that logs request method and URL.
// Request/response bodies are NOT logged to prevent sensitive data exposure.
// Use LogPlusMiddlewareWithConfig for full control over logging behavior.
func LogPlusMiddleware() gin.HandlerFunc {
	return LogPlusMiddlewareWithConfig(LogPlusConfig{
		LogBody:    false,
		LogHeaders: true,
	})
}

// LogPlusMiddlewareWithConfig returns a Gin middleware with configurable logging behavior.
// WARNING: Setting LogBody to true may expose sensitive data (passwords, tokens, PII) in logs.
func LogPlusMiddlewareWithConfig(cfg LogPlusConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r reqLog
		r.ReqID = utils.RandomString()
		r.Method = c.Request.Method
		r.URL = c.Request.URL.RequestURI()
		r.In = time.Now()

		if cfg.LogHeaders {
			r.Header = redactHeaders(c.Request.Header)
		}

		rawReqData, _ := io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawReqData))
		if cfg.LogBody {
			r.Body = string(rawReqData)
		} else {
			r.Body = "[REDACTED]"
		}

		customWriter := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = customWriter
		c.Next()

		r.Out = time.Now()
		if cfg.LogBody {
			r.Resp = customWriter.body.String()
		} else {
			r.Resp = "[REDACTED]"
		}

		slog.Debug("LPM", slog.Any("info", r))
		slog.Info("LPM", slog.String("req-id", r.ReqID), slog.String("catch", r.Method+" "+r.URL))
	}
}
