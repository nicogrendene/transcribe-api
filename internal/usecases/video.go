package usecases

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/pkg/utils"
)

// VideoUseCaseImpl implementa la lógica de videos
type VideoUseCaseImpl struct {
	config config.Config
}

// NewVideoUseCase crea una nueva instancia del use case de videos
func NewVideoUseCase(config config.Config) VideoUseCase {
	return &VideoUseCaseImpl{
		config: config,
	}
}

// GetVideos obtiene la lista de videos
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

// GetVideo obtiene la ruta de un video
func (v *VideoUseCaseImpl) GetVideo(filename string) (string, error) {
	if !utils.ValidateFilename(filename) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	videoPath := filepath.Join(v.config.VideosPath, filename)

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de video no encontrado")
	}

	return videoPath, nil
}

// GetSubtitles obtiene la ruta de los subtítulos
func (v *VideoUseCaseImpl) GetSubtitles(filename string) (string, error) {
	if !utils.ValidateFilename(filename) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	// Convertir .mp4 a .vtt
	baseName := strings.TrimSuffix(filename, ".mp4")
	subtitleFilename := baseName + ".vtt"
	subtitlePath := filepath.Join(v.config.VideosPath, subtitleFilename)

	if _, err := os.Stat(subtitlePath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de subtítulos no encontrado")
	}

	return subtitlePath, nil
}

// GetThumbnail obtiene la ruta de la miniatura
func (v *VideoUseCaseImpl) GetThumbnail(filename string) (string, error) {
	if !utils.ValidateFilename(filename) {
		return "", fmt.Errorf("nombre de archivo inválido")
	}

	// Convertir .mp4 a .jpg
	baseName := strings.TrimSuffix(filename, ".mp4")
	thumbnailFilename := baseName + ".jpg"
	thumbnailPath := filepath.Join(v.config.VideosPath, thumbnailFilename)

	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo de miniatura no encontrado")
	}

	return thumbnailPath, nil
}
