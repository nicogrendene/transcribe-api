package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
)

type OpenAIService struct {
	APIKey         string
	Model          string
	PricePer1K     float64
	ChatModel      string
	ChatPricePer1K float64
}

// NewOpenAIService crea una nueva instancia del servicio OpenAI
func NewOpenAIService(apiKey, model string, pricePer1K float64, chatModel string, chatPricePer1K float64) (*OpenAIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if model == "" {
		return nil, fmt.Errorf("model is required")
	}
	if pricePer1K <= 0 {
		return nil, fmt.Errorf("price per 1K is required")
	}
	if chatModel == "" {
		return nil, fmt.Errorf("chat model is required")
	}
	if chatPricePer1K <= 0 {
		return nil, fmt.Errorf("chat price per 1K is required")
	}

	return &OpenAIService{
		APIKey:         apiKey,
		Model:          model,
		PricePer1K:     pricePer1K,
		ChatModel:      chatModel,
		ChatPricePer1K: chatPricePer1K,
	}, nil
}

// GenerateEmbedding genera un embedding para el texto dado
func (s *OpenAIService) GenerateEmbedding(ctx context.Context, text string) ([]float32, int, error) {
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
	log.Info(ctx, "Embedding generated",
		log.Any("tokens", tokens),
		log.Float("costo", costo),
	)

	return embResp.Data[0].Embedding, tokens, nil
}

// GenerateAnswer genera una respuesta usando el modelo de chat de OpenAI
func (s *OpenAIService) GenerateAnswer(ctx context.Context, query string, contextTexts []string) (string, int, error) {
	searchContext := ""
	for i, text := range contextTexts {
		searchContext += fmt.Sprintf("Fragmento %d: %s\n\n", i+1, text)
	}

	systemPrompt := `Eres un asistente que responde preguntas basándote únicamente en el contexto proporcionado. 
Responde de manera clara y concisa. Si la información no está disponible en el contexto, 
indica que no tienes suficiente información para responder la pregunta.

Contexto:
` + searchContext

	messages := []models.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: query,
		},
	}

	reqBody := models.OpenAIChatRequest{
		Model:       s.ChatModel,
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error calling OpenAI API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusTooManyRequests {
			return "", 0, fmt.Errorf("cuota de OpenAI excedida")
		}
		return "", 0, fmt.Errorf("OpenAI API retornó status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp models.OpenAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", 0, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", 0, fmt.Errorf("no se recibió respuesta de OpenAI")
	}

	tokens := chatResp.Usage.TotalTokens
	costo := float64(tokens) * s.ChatPricePer1K / 1000.0
	log.Info(ctx, "Got answer", log.Any("tokens", tokens), log.Float("costo", costo))

	return chatResp.Choices[0].Message.Content, tokens, nil
}
