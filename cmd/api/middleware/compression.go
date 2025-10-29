package middleware

import (
	"compress/gzip"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// GzipMiddleware proporciona compresión gzip para respuestas
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo comprimir si el cliente acepta gzip y el contenido es apropiado
		if !shouldCompress(c) {
			c.Next()
			return
		}

		// Crear un writer gzip
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		// Configurar headers de compresión
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Reemplazar el writer de la respuesta
		c.Writer = &gzipResponseWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		c.Next()
	}
}

// shouldCompress determina si la respuesta debe ser comprimida
func shouldCompress(c *gin.Context) bool {
	// Verificar si el cliente acepta gzip
	acceptEncoding := c.GetHeader("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		return false
	}

	// Solo comprimir ciertos tipos de contenido
	contentType := c.GetHeader("Content-Type")
	switch {
	case strings.HasPrefix(contentType, "video/"):
		// Los videos ya están comprimidos, no comprimir de nuevo
		return false
	case strings.HasPrefix(contentType, "image/"):
		// Las imágenes ya están comprimidas, no comprimir de nuevo
		return false
	case strings.HasPrefix(contentType, "application/json"):
		// Comprimir JSON
		return true
	case strings.HasPrefix(contentType, "text/"):
		// Comprimir texto
		return true
	default:
		// No comprimir otros tipos
		return false
	}
}

// gzipResponseWriter envuelve el ResponseWriter para compresión gzip
type gzipResponseWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipResponseWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}
