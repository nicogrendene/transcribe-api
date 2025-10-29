package usecases

import (
	"context"

	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
)

type SearchUseCase interface {
	Search(ctx context.Context, query string, topK int) (*models.SearchResponse, error)
}

type HealthUseCase interface {
	CheckHealth(ctx context.Context) (*models.HealthResponse, error)
}

type StatsUseCase interface {
	GetStats(ctx context.Context) (*models.StatsResponse, error)
}

type VideoUseCase interface {
	GetVideos(ctx context.Context) ([]byte, error)
	GetVideo(ctx context.Context, id string) (string, error)
	GetVideoWithQuality(ctx context.Context, id string, quality string) (string, error)
	GetSubtitles(ctx context.Context, id string) (string, error)
	GetThumbnail(ctx context.Context, id string) (string, error)
	GetSummary(ctx context.Context, id string) (string, error)
}
