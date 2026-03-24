package apimiddleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/johnpoint/go-bootstrap/v2/utils"
	"io"
	"log/slog"
	"net/http"
	"time"
)

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

// LogPlusMiddleware returns a Gin middleware that logs detailed request and response information.
// It captures request headers, body, URL, method, and response data with a unique request ID.
func LogPlusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r reqLog
		r.ReqID = utils.RandomString()
		r.Header = c.Request.Header
		r.Method = c.Request.Method
		rawReqData, _ := io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawReqData))

		r.Body = string(rawReqData)
		r.URL = c.Request.URL.RequestURI()
		customWriter := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = customWriter
		// 处理请求
		c.Next()

		r.Resp = customWriter.body.String()
		r.Out = time.Now()
		slog.Debug("LPM", slog.Any("info", r))
		slog.Info("LPM", slog.String("catch", r.Method+" "+r.URL))
	}
}
