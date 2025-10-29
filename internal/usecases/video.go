package usecases

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/pkg/utils"
)

type VideoUseCaseImpl struct {
	config config.Config
}

func NewVideoUseCase(config config.Config) VideoUseCase {
	return &VideoUseCaseImpl{
		config: config,
	}
}

func (v *VideoUseCaseImpl) GetVideos(ctx context.Context) ([]byte, error) {
	jsonFile, err := os.Open("videos.json")
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo de videos: %v", err)
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("no se pudo procesar el archivo de videos: %v", err)
	}

	return jsonData, nil
}

func (v *VideoUseCaseImpl) GetVideo(ctx context.Context, filename string) (string, error) {
	return v.GetVideoWithQuality(ctx, filename, "")
}

func (v *VideoUseCaseImpl) GetVideoWithQuality(ctx context.Context, filename string, quality string) (string, error) {
	if !utils.ValidateFilename(filename) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	// Determinar el nombre del archivo basado en la calidad
	var videoFile string
	switch quality {
	case "high", "1080p":
		videoFile = "video_1080p.mp4"
	case "medium", "720p":
		videoFile = "video_720p.mp4"
	case "low", "480p":
		videoFile = "video_480p.mp4"
	case "mobile", "360p":
		videoFile = "video_360p.mp4"
	default:
		// Calidad por defecto
		videoFile = "video.mp4"
	}

	videoPath := filepath.Join(v.config.VideosPath, filename, videoFile)

	// Si el archivo de calidad específica no existe, intentar con el archivo por defecto
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		if quality != "" {
			// Intentar con el archivo por defecto
			videoPath = filepath.Join(v.config.VideosPath, filename, "video.mp4")
			if _, err := os.Stat(videoPath); os.IsNotExist(err) {
				return "", fmt.Errorf("archivo de video no encontrado")
			}
		} else {
			return "", fmt.Errorf("archivo de video no encontrado")
		}
	}

	return videoPath, nil
}

func (v *VideoUseCaseImpl) GetSubtitles(ctx context.Context, id string) (string, error) {
	if !utils.ValidateFilename(id) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	subtitlePath := filepath.Join(v.config.VideosPath, id, "subtitles.vtt")

	if _, err := os.Stat(subtitlePath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de subtítulos no encontrado")
	}

	return subtitlePath, nil
}

func (v *VideoUseCaseImpl) GetThumbnail(ctx context.Context, id string) (string, error) {
	if !utils.ValidateFilename(id) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	thumbnailPath := filepath.Join(v.config.VideosPath, id, "thumbnail.jpg")

	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de miniatura no encontrado")
	}

	return thumbnailPath, nil
}

func (v *VideoUseCaseImpl) GetSummary(ctx context.Context, id string) (string, error) {
	if !utils.ValidateFilename(id) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	summaryPath := filepath.Join(v.config.VideosPath, id, "summary.txt")

	file, err := os.Open(summaryPath)
	if err != nil {
		return "", fmt.Errorf("no se encontrado por archivo de summary.txt")
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("no se encontrado por archivo de summary.txt")
	}

	return string(content), nil
}
