package usecases

import (
	"context"
	"fmt"

	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/services"
)

// StatsUseCaseImpl implementa la lógica de estadísticas
type StatsUseCaseImpl struct {
	pineconeService *services.PineconeService
}

// NewStatsUseCase crea una nueva instancia del use case de stats
func NewStatsUseCase(pineconeService *services.PineconeService) StatsUseCase {
	return &StatsUseCaseImpl{
		pineconeService: pineconeService,
	}
}

// GetStats obtiene las estadísticas del sistema
func (s *StatsUseCaseImpl) GetStats(ctx context.Context) (*models.StatsResponse, error) {
	stats, err := s.pineconeService.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo estadísticas: %v", err)
	}
	return stats, nil
}
