package usecases

import "github.com/ngrendenebos/scripts/transcribe-api/internal/models"

type SearchUseCase interface {
	Search(query string, topK int) (*models.BusquedaResponse, error)
}

type HealthUseCase interface {
	CheckHealth() (*models.HealthResponse, error)
}

type StatsUseCase interface {
	GetStats() (*models.StatsResponse, error)
}

type VideoUseCase interface {
	GetVideos() ([]byte, error)
	GetVideo(id string) (string, error)
	GetSubtitles(id string) (string, error)
	GetThumbnail(id string) (string, error)
	GetSummary(id string) (string, error)
}
