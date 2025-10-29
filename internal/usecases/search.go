package usecases

import (
	"context"
	"fmt"

	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
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

// Search realiza una búsqueda vectorial y genera una respuesta con OpenAI
func (s *SearchUseCaseImpl) Search(ctx context.Context, query string, topK int) (*models.SearchResponse, error) {
	// Validar parámetros
	if query == "" {
		return nil, fmt.Errorf("query no puede estar vacío")
	}

	if topK < 1 || topK > s.config.MaxTopK {
		return nil, fmt.Errorf("top_k debe estar entre 1 y %d", s.config.MaxTopK)
	}

	// Generar embedding
	embedding, tokens, err := s.openaiService.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error generando embedding: %v", err)
	}

	// Calcular costo inicial
	costo := float64(tokens) * s.config.EmbeddingPricePer1K / 1000.0

	// Buscar en Pinecone
	res, err := s.pineconeService.Search(ctx, embedding, topK)
	if err != nil {
		return nil, fmt.Errorf("error en búsqueda: %v", err)
	}

	filtrados := s.filterByScore(res, s.config.MinScoreThreshold)

	// Generar respuesta con OpenAI si hay resultados
	var generatedAnswer string
	if len(filtrados) > 0 {
		// Extraer textos de los resultados
		contextTexts := make([]string, 0, len(filtrados))
		for _, result := range filtrados {
			if result.Text != "" {
				contextTexts = append(contextTexts, result.Text)
			}
		}

		// Generar respuesta con OpenAI
		answer, chatTokens, err := s.openaiService.GenerateAnswer(ctx, query, contextTexts)
		if err != nil {
			log.Error(ctx, "Error generando respuesta", log.Err(err))
		} else {
			generatedAnswer = answer
			// Agregar costo de chat
			costo += float64(chatTokens) * s.config.ChatPricePer1K / 1000.0
		}
	}

	return &models.SearchResponse{
		Query:           query,
		Results:         filtrados,
		Total:           len(filtrados),
		GeneratedAnswer: generatedAnswer,
		CostoUSD:        costo,
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
