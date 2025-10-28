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
func (s *SearchUseCaseImpl) Search(query string, topK int) (*models.SearchResponse, error) {
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
	res, err := s.pineconeService.Search(embedding, topK)
	if err != nil {
		return nil, fmt.Errorf("error en búsqueda: %v", err)
	}

	filtrados := s.filterByScore(res, s.config.MinScoreThreshold)

	return &models.SearchResponse{
		Query:    query,
		Results:  filtrados,
		Total:    len(filtrados),
		CostoUSD: costo,
	}, nil
}

// filterByScore filtra resultados por umbral de similitud
func (s *SearchUseCaseImpl) filterByScore(resultados []models.ChunkResponse, threshold float64) []models.ChunkResponse {
	filtrados := make([]models.ChunkResponse, 0, len(resultados))
	for _, r := range resultados {
		if float64(r.Score) >= threshold {
			filtrados = append(filtrados, r)
		}
	}
	return filtrados
}
