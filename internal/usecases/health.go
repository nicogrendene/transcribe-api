package usecases

import (
	"context"

	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/services"
)

type HealthUseCaseImpl struct {
	pineconeService *services.PineconeService
}

func NewHealthUseCase(pineconeService *services.PineconeService) HealthUseCase {
	return &HealthUseCaseImpl{
		pineconeService: pineconeService,
	}
}

func (h *HealthUseCaseImpl) CheckHealth(ctx context.Context) (*models.HealthResponse, error) {
	if h.pineconeService == nil {
		return &models.HealthResponse{
			Status:  "error",
			Message: "Servicio no inicializado",
		}, nil
	}

	return &models.HealthResponse{
		Status:  "healthy",
		Message: "API funcionando correctamente",
	}, nil
}
