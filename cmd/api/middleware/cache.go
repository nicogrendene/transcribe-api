package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// VideoCacheMiddleware proporciona cache agresivo para videos
func VideoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo aplicar cache a rutas de video
		if !isVideoRoute(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Headers de cache agresivo para videos
		c.Header("Cache-Control", "public, max-age=31536000, immutable") // 1 año
		c.Header("Expires", time.Now().Add(365*24*time.Hour).Format(http.TimeFormat))
		c.Header("Vary", "Accept-Encoding, Range")

		// Headers para optimización de video
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")

		// Permitir preload de recursos
		c.Header("Link", "</video>; rel=preload; as=video")

		c.Next()
	}
}

// ThumbnailCacheMiddleware proporciona cache para thumbnails
func ThumbnailCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo aplicar cache a rutas de thumbnail
		if !isThumbnailRoute(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Headers de cache para thumbnails
		c.Header("Cache-Control", "public, max-age=604800, immutable") // 1 semana
		c.Header("Expires", time.Now().Add(7*24*time.Hour).Format(http.TimeFormat))
		c.Header("Vary", "Accept-Encoding")

		// Headers para optimización de imagen
		c.Header("X-Content-Type-Options", "nosniff")

		c.Next()
	}
}

// isVideoRoute verifica si la ruta es para servir videos
func isVideoRoute(path string) bool {
	return path == "/video" ||
		(len(path) > 7 && path[:7] == "/video/") ||
		(len(path) > 7 && path[:7] == "/video/")
}

// isThumbnailRoute verifica si la ruta es para servir thumbnails
func isThumbnailRoute(path string) bool {
	return len(path) > 12 && path[len(path)-10:] == "/thumbnail"
}
