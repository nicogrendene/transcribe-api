package usecases

import (
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

func (v *VideoUseCaseImpl) GetVideos() ([]byte, error) {
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

func (v *VideoUseCaseImpl) GetVideo(filename string) (string, error) {
	if !utils.ValidateFilename(filename) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	videoPath := filepath.Join(v.config.VideosPath, filename, "video.mp4")

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de video no encontrado")
	}

	return videoPath, nil
}

func (v *VideoUseCaseImpl) GetSubtitles(id string) (string, error) {
	if !utils.ValidateFilename(id) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	subtitlePath := filepath.Join(v.config.VideosPath, id, "subtitles.vtt")

	if _, err := os.Stat(subtitlePath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de subtítulos no encontrado")
	}

	return subtitlePath, nil
}

func (v *VideoUseCaseImpl) GetThumbnail(id string) (string, error) {
	if !utils.ValidateFilename(id) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	thumbnailPath := filepath.Join(v.config.VideosPath, id, "thumbnail.jpg")

	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de miniatura no encontrado")
	}

	return thumbnailPath, nil
}

func (v *VideoUseCaseImpl) GetSummary(id string) (string, error) {
	if !utils.ValidateFilename(id) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	summaryPath := filepath.Join(v.config.VideosPath, id, "summary.txt")

	if _, err := os.Stat(summaryPath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de resumen no encontrado")
	}

	return summaryPath, nil
}
