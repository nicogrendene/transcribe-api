package middleware

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
)

type RequestIDKey struct{}

type responseWriterWrapper struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriterWrapper) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()

		start := time.Now()

		ctx := log.With(c.Request.Context(),
			log.String("request_id", requestID),
			log.String("method", c.Request.Method),
			log.String("path", c.Request.URL.Path),
			log.String("query", c.Request.URL.RawQuery),
			log.String("remote_addr", c.ClientIP()),
			log.String("user_agent", c.Request.UserAgent()),
		)

		ctx = context.WithValue(ctx, RequestIDKey{}, requestID)

		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-ID", requestID)

		blw := &responseWriterWrapper{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		status := c.Writer.Status()
		responseBody := blw.body.String()
		if len(responseBody) > 1000 {
			responseBody = responseBody[:1000] + "..."
		}
		fields := []log.Field{
			log.Int("status", status),
			log.Duration("duration", time.Since(start)),
			log.String("response_body", responseBody),
		}

		switch {
		case status >= 500:
			log.Error(ctx, "Request processed", fields...)
		case status >= 400:
			log.Warn(ctx, "Request processed", fields...)
		default:
			log.Info(ctx, "Request processed", fields...)
		}
	}
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return requestID
	}
	return ""
}
