package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
)

// OpenAIService maneja las interacciones con OpenAI
type OpenAIService struct {
	APIKey     string
	Model      string
	PricePer1K float64
}

// NewOpenAIService crea una nueva instancia del servicio OpenAI
func NewOpenAIService(apiKey, model string, pricePer1K float64) (*OpenAIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if model == "" {
		return nil, fmt.Errorf("model is required")
	}
	if pricePer1K <= 0 {
		return nil, fmt.Errorf("price per 1K is required")
	}

	return &OpenAIService{
		APIKey:     apiKey,
		Model:      model,
		PricePer1K: pricePer1K,
	}, nil
}

// GenerateEmbedding genera un embedding para el texto dado
func (s *OpenAIService) GenerateEmbedding(text string) ([]float32, int, error) {
	reqBody := models.OpenAIEmbeddingRequest{
		Input:      text,
		Model:      s.Model,
		Dimensions: 512, // Usar constante del config
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error calling OpenAI API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, 0, fmt.Errorf("cuota de OpenAI excedida. Verifica tu plan en https://platform.openai.com/account/billing")
		}
		return nil, 0, fmt.Errorf("OpenAI API retornó status %d: %s", resp.StatusCode, string(body))
	}

	var embResp models.OpenAIEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, 0, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if len(embResp.Data) == 0 {
		return nil, 0, fmt.Errorf("no se recibieron embeddings")
	}

	tokens := embResp.Usage.TotalTokens
	costo := float64(tokens) * s.PricePer1K / 1000.0
	log.Printf("💰 Embedding: %d tokens, $%.6f", tokens, costo)

	return embResp.Data[0].Embedding, tokens, nil
}
