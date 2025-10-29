package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

func ServeVideo(videoUseCase usecases.VideoUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := log.With(c.Request.Context(), log.UseCase("serve_video"))

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "id parameter is required",
			})
			return
		}

		// Obtener parámetro de calidad de la query string
		quality := c.Query("quality")

		videoPath, err := videoUseCase.GetVideoWithQuality(ctx, id, quality)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Video file not found",
				Details: err.Error(),
			})
			return
		}

		// Abrir el archivo de video
		file, err := os.Open(videoPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error opening video file",
				Details: err.Error(),
			})
			return
		}
		defer file.Close()

		// Obtener información del archivo
		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error getting file info",
				Details: err.Error(),
			})
			return
		}

		fileSize := fileInfo.Size()

		// Headers básicos
		c.Header("Content-Type", "video/mp4")
		c.Header("Accept-Ranges", "bytes")
		c.Header("Cache-Control", "public, max-age=86400") // 24 horas de cache
		c.Header("ETag", fmt.Sprintf("\"%x-%x\"", fileInfo.ModTime().Unix(), fileSize))
		c.Header("Last-Modified", fileInfo.ModTime().Format(http.TimeFormat))

		// Headers para preloading y optimización
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")

		// Headers de preload para recursos relacionados
		baseURL := c.Request.URL.Scheme + "://" + c.Request.Host
		preloadLinks := []string{
			fmt.Sprintf("<%s/video/%s/thumbnail>; rel=preload; as=image", baseURL, id),
			fmt.Sprintf("<%s/video/%s/subtitles>; rel=preload; as=fetch", baseURL, id),
			fmt.Sprintf("<%s/video/%s/summary>; rel=preload; as=fetch", baseURL, id),
		}
		c.Header("Link", strings.Join(preloadLinks, ", "))

		// Verificar si el cliente solicita un rango específico
		rangeHeader := c.GetHeader("Range")
		if rangeHeader == "" {
			// Sin Range header, servir el archivo completo
			c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
			c.Status(http.StatusOK)
			io.Copy(c.Writer, file)
			return
		}

		// Parsear el Range header
		ranges, err := parseRange(rangeHeader, fileSize)
		if err != nil {
			c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// Si hay múltiples rangos, usar el primero (soporte básico)
		start := ranges[0].start
		end := ranges[0].end

		// Asegurar que el rango sea válido
		if start >= fileSize || end >= fileSize || start > end {
			c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// Calcular el tamaño del contenido
		contentLength := end - start + 1

		// Headers para respuesta parcial
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
		c.Status(http.StatusPartialContent)

		// Posicionar el archivo en el inicio del rango
		file.Seek(start, 0)

		// Crear un reader limitado para el rango solicitado
		limitedReader := io.LimitReader(file, contentLength)
		io.Copy(c.Writer, limitedReader)
	}
}

// Estructura para representar un rango de bytes
type httpRange struct {
	start, end int64
}

// parseRange parsea el header Range y retorna los rangos solicitados
func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, fmt.Errorf("header not present")
	}

	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, fmt.Errorf("invalid range")
	}

	var ranges []httpRange
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}

		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, fmt.Errorf("invalid range")
		}

		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r httpRange

		if start == "" {
			// Si no hay start, significa "últimos N bytes"
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil || i < 0 {
				return nil, fmt.Errorf("invalid range")
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.end = size - 1
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i < 0 {
				return nil, fmt.Errorf("invalid range")
			}
			if i >= size {
				return nil, fmt.Errorf("invalid range")
			}
			r.start = i

			if end == "" {
				// Si no hay end, significa "desde start hasta el final"
				r.end = size - 1
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					return nil, fmt.Errorf("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				r.end = i
			}
		}

		ranges = append(ranges, r)
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("no ranges found")
	}

	return ranges, nil
}
