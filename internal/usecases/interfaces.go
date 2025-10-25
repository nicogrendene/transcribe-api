package usecases

import "github.com/ngrendenebos/scripts/transcribe-api/internal/models"

// Interfaces para los use cases

// SearchUseCase define la lógica de negocio para búsquedas
type SearchUseCase interface {
	Search(query string, topK int) (*models.BusquedaResponse, error)
}

// HealthUseCase define la lógica de negocio para health checks
type HealthUseCase interface {
	CheckHealth() (*models.HealthResponse, error)
}

// StatsUseCase define la lógica de negocio para estadísticas
type StatsUseCase interface {
	GetStats() (*models.StatsResponse, error)
}

// VideoUseCase define la lógica de negocio para videos
type VideoUseCase interface {
	GetVideos() ([]byte, error)
	GetVideo(filename string) (string, error)
	GetSubtitles(filename string) (string, error)
	GetThumbnail(filename string) (string, error)
}
