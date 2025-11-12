package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// ServeThumbnail retorna un handler para servir miniaturas
func ServeThumbnail(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := log.With(c.Request.Context(), log.UseCase("serve_thumbnail"))

		id := c.Param("id")

		thumbnailPath, err := videoUseCase.GetThumbnail(ctx, id)
		if err != nil {
			// Si no existe, devolver un placeholder SVG simple
			svgPlaceholder := `<svg width="320" height="180" xmlns="http://www.w3.org/2000/svg">
				<rect width="320" height="180" fill="#374151"/>
				<text x="160" y="90" text-anchor="middle" fill="white" font-family="Arial" font-size="24">🎥</text>
			</svg>`

			c.Header("Content-Type", "image/svg+xml")
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.Header("Access-Control-Allow-Origin", "*")
			c.String(http.StatusOK, svgPlaceholder)
			return
		}

		c.Header("Content-Type", "image/jpeg")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("Access-Control-Allow-Origin", "*")
		c.File(thumbnailPath)
	}
}
