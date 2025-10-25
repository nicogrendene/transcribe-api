package usecases

import (
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/services"
)

// HealthUseCaseImpl implementa la lógica de health check
type HealthUseCaseImpl struct {
	pineconeService *services.PineconeService
}

// NewHealthUseCase crea una nueva instancia del use case de health
func NewHealthUseCase(pineconeService *services.PineconeService) HealthUseCase {
	return &HealthUseCaseImpl{
		pineconeService: pineconeService,
	}
}

// CheckHealth verifica el estado de salud del sistema
func (h *HealthUseCaseImpl) CheckHealth() (*models.HealthResponse, error) {
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
