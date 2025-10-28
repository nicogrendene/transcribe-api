package dependencies

import (
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/services"
)

type Dependencies struct {
	PineconeService *services.PineconeService
	OpenAIService   *services.OpenAIService
}

func NewDependencies(cfg config.Config) (Dependencies, error) {
	var deps Dependencies

	openAIService, err := services.NewOpenAIService(
		cfg.OpenAIAPIKey,
		cfg.EmbeddingModel,
		cfg.EmbeddingPricePer1K,
		cfg.ChatModel,
		cfg.ChatPricePer1K,
	)
	if err != nil {
		return deps, err
	}
	deps.OpenAIService = openAIService

	pineconeService, err := services.NewPineconeService(
		cfg.PineconeAPIKey,
		cfg.IndexName,
	)
	if err != nil {
		return deps, err
	}
	deps.PineconeService = pineconeService

	return deps, nil
}
