package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	// API Keys
	OpenAIAPIKey   string
	PineconeAPIKey string

	// Pinecone
	IndexName          string
	EmbeddingModel     string
	EmbeddingDimension int

	// Umbrales y límites
	MinScoreThreshold float64
	MaxTopK           int
	DefaultTopK       int

	// Precios
	EmbeddingPricePer1K float64

	// Servidor
	Port string

	// Rutas
	VideosPath string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() (Config, error) {
	var config Config

	if err := godotenv.Load(); err != nil {
		return config, fmt.Errorf("no se pudo cargar archivo .env: %v", err)
	}

	// Variables opcionales (con valores por defecto)
	config.OpenAIAPIKey = getEnvOrDefault("OPENAI_API_KEY", "")
	config.PineconeAPIKey = getEnvOrDefault("PINECONE_API_KEY", "")
	config.IndexName = getEnvOrDefault("INDEX_NAME", config.IndexName)
	config.Port = getEnvOrDefault("PORT", config.Port)
	config.VideosPath = getEnvOrDefault("VIDEOS_PATH", config.VideosPath)

	// Variables numéricas opcionales
	if threshold := getEnvOrDefault("MIN_SCORE_THRESHOLD", ""); threshold != "" {
		if f, err := strconv.ParseFloat(threshold, 64); err == nil {
			config.MinScoreThreshold = f
		}
	}

	if err := config.validate(); err != nil {
		return config, err
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) validate() error {
	if c.OpenAIAPIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY es requerida")
	}

	if c.PineconeAPIKey == "" {
		return fmt.Errorf("PINECONE_API_KEY es requerida")
	}

	if c.MinScoreThreshold < 0 || c.MinScoreThreshold > 1 {
		return fmt.Errorf("MIN_SCORE_THRESHOLD debe estar entre 0 y 1")
	}

	if c.MaxTopK < 1 || c.DefaultTopK < 1 {
		return fmt.Errorf("top_k debe ser mayor a 0")
	}

	return nil
}
