package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
)

// RecoveryWithLogging middleware que captura panics y los loggea apropiadamente
func RecoveryWithLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Verificar si es una conexión rota
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Obtener contexto con request ID si está disponible
				ctx := c.Request.Context()
				requestID := GetRequestID(ctx)

				// Preparar campos de log
				fields := []log.Field{
					log.String("request_id", requestID),
					log.String("method", c.Request.Method),
					log.String("path", c.Request.URL.Path),
					log.String("error", fmt.Sprintf("%v", err)),
					log.String("stack", string(debug.Stack())),
				}

				// Agregar información adicional del request
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				fields = append(fields, log.String("request", string(httpRequest)))

				// Loggear el panic
				if brokenPipe {
					log.Error(ctx, "Broken pipe error", fields...)
				} else {
					log.Panic(ctx, "Panic recovered", fields...)
				}

				// Si la conexión está rota, no podemos enviar respuesta
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				// Responder con error 500
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// LoggingWithMetrics middleware que incluye métricas de performance
func LoggingWithMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Procesar request
		c.Next()

		// Calcular métricas
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		// Obtener contexto
		ctx := c.Request.Context()
		requestID := GetRequestID(ctx)

		// Preparar campos de log
		fields := []log.Field{
			log.String("request_id", requestID),
			log.String("method", method),
			log.String("path", path),
			log.String("query", raw),
			log.String("client_ip", clientIP),
			log.Int("status", statusCode),
			log.Duration("latency", latency),
			log.Int("body_size", bodySize),
			log.String("user_agent", c.Request.UserAgent()),
		}

		// Loggear con nivel apropiado
		switch {
		case statusCode >= 500:
			log.Error(ctx, "HTTP Request", fields...)
		case statusCode >= 400:
			log.Warn(ctx, "HTTP Request", fields...)
		default:
			log.Info(ctx, "HTTP Request", fields...)
		}
	}
}
