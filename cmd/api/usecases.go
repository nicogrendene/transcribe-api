package main

import (
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/dependencies"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// Usecases contiene todos los use cases de la aplicación
type Usecases struct {
	SearchUseCase usecases.SearchUseCase
	HealthUseCase usecases.HealthUseCase
	StatsUseCase  usecases.StatsUseCase
	VideoUseCase  usecases.VideoUseCase
}

// NewUsecases crea una nueva instancia de use cases
func NewUsecases(deps dependencies.Dependencies, cfg config.Config) Usecases {
	return Usecases{
		SearchUseCase: usecases.NewSearchUseCase(deps.OpenAIService, deps.PineconeService, cfg),
		HealthUseCase: usecases.NewHealthUseCase(deps.PineconeService),
		StatsUseCase:  usecases.NewStatsUseCase(deps.PineconeService),
		VideoUseCase:  usecases.NewVideoUseCase(cfg),
	}
}
