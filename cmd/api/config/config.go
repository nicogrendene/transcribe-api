package config

import (
	"fmt"
	"log"
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

	// Variables requeridas
	config.OpenAIAPIKey = getEnvOrDefault("OPENAI_API_KEY", "")
	config.PineconeAPIKey = getEnvOrDefault("PINECONE_API_KEY", "")
	config.IndexName = getEnvOrDefault("INDEX_NAME", "")
	config.EmbeddingModel = getEnvOrDefault("EMBEDDING_MODEL", "")
	config.Port = getEnvOrDefault("PORT", "")
	config.VideosPath = getEnvOrDefault("VIDEOS_PATH", "")

	// Variables numéricas opcionales
	if threshold := getEnvOrDefault("MIN_SCORE_THRESHOLD", ""); threshold != "" {
		if f, err := strconv.ParseFloat(threshold, 64); err == nil {
			config.MinScoreThreshold = f
		}
	}

	if dimension := getEnvOrDefault("EMBEDDING_DIMENSION", ""); dimension != "" {
		if i, err := strconv.Atoi(dimension); err == nil {
			config.EmbeddingDimension = i
		}
	}

	if maxTopK := getEnvOrDefault("MAX_TOP_K", ""); maxTopK != "" {
		if i, err := strconv.Atoi(maxTopK); err == nil {
			config.MaxTopK = i
		}
	}

	if defaultTopK := getEnvOrDefault("DEFAULT_TOP_K", ""); defaultTopK != "" {
		if i, err := strconv.Atoi(defaultTopK); err == nil {
			config.DefaultTopK = i
		}
	}

	if price := getEnvOrDefault("EMBEDDING_PRICE_PER_1K", ""); price != "" {
		if f, err := strconv.ParseFloat(price, 64); err == nil {
			config.EmbeddingPricePer1K = f
		}
	}

	if err := config.validate(); err != nil {
		return config, err
	}

	log.Println("✅ Configuración cargada correctamente")
	log.Println("OPENAI_API_KEY: " + config.OpenAIAPIKey)
	log.Println("PINECONE_API_KEY: " + config.PineconeAPIKey)
	log.Println("INDEX_NAME: " + config.IndexName)
	log.Println("EMBEDDING_MODEL: " + config.EmbeddingModel)
	log.Println("EMBEDDING_DIMENSION: " + strconv.Itoa(config.EmbeddingDimension))
	log.Println("MIN_SCORE_THRESHOLD: " + strconv.FormatFloat(config.MinScoreThreshold, 'f', -1, 64))
	log.Println("MAX_TOP_K: " + strconv.Itoa(config.MaxTopK))
	log.Println("DEFAULT_TOP_K: " + strconv.Itoa(config.DefaultTopK))
	log.Println("EMBEDDING_PRICE_PER_1K: " + strconv.FormatFloat(config.EmbeddingPricePer1K, 'f', -1, 64))
	log.Println("PORT: " + config.Port)
	log.Println("VIDEOS_PATH: " + config.VideosPath)

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

	if c.IndexName == "" {
		return fmt.Errorf("INDEX_NAME es requerida")
	}

	if c.Port == "" {
		return fmt.Errorf("PORT es requerida")
	}

	if c.VideosPath == "" {
		return fmt.Errorf("VIDEOS_PATH es requerida")
	}

	if c.MinScoreThreshold < 0 || c.MinScoreThreshold > 1 {
		return fmt.Errorf("MIN_SCORE_THRESHOLD debe estar entre 0 y 1")
	}

	if c.MaxTopK < 1 || c.DefaultTopK < 1 {
		return fmt.Errorf("top_k debe ser mayor a 0")
	}

	return nil
}
