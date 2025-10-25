package usecases

import (
	"fmt"

	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/services"
)

// SearchUseCaseImpl implementa la lógica de búsqueda
type SearchUseCaseImpl struct {
	openaiService   *services.OpenAIService
	pineconeService *services.PineconeService
	config          config.Config
}

// NewSearchUseCase crea una nueva instancia del use case de búsqueda
func NewSearchUseCase(openaiService *services.OpenAIService, pineconeService *services.PineconeService, config config.Config) SearchUseCase {
	return &SearchUseCaseImpl{
		openaiService:   openaiService,
		pineconeService: pineconeService,
		config:          config,
	}
}

// Search realiza una búsqueda vectorial
func (s *SearchUseCaseImpl) Search(query string, topK int) (*models.BusquedaResponse, error) {
	// Validar parámetros
	if query == "" {
		return nil, fmt.Errorf("query no puede estar vacío")
	}

	if topK < 1 || topK > s.config.MaxTopK {
		return nil, fmt.Errorf("top_k debe estar entre 1 y %d", s.config.MaxTopK)
	}

	// Generar embedding
	embedding, tokens, err := s.openaiService.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("error generando embedding: %v", err)
	}

	// Calcular costo
	costo := float64(tokens) * s.config.EmbeddingPricePer1K / 1000.0

	// Buscar en Pinecone
	resultados, err := s.pineconeService.Search(embedding, topK)
	if err != nil {
		return nil, fmt.Errorf("error en búsqueda: %v", err)
	}

	// Filtrar por umbral de similitud
	filtrados := s.filterByScore(resultados, s.config.MinScoreThreshold)

	return &models.BusquedaResponse{
		Query:      query,
		Resultados: filtrados,
		Total:      len(filtrados),
		CostoUSD:   costo,
	}, nil
}

// filterByScore filtra resultados por umbral de similitud
func (s *SearchUseCaseImpl) filterByScore(resultados []models.ChunkResultado, threshold float64) []models.ChunkResultado {
	filtrados := make([]models.ChunkResultado, 0, len(resultados))
	for _, r := range resultados {
		if float64(r.Score) >= threshold {
			filtrados = append(filtrados, r)
		}
	}
	return filtrados
}
