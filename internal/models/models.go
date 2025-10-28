package models

type SearchRequest struct {
	Query string `json:"query" binding:"required,min=2"`
	TopK  int    `json:"top_k"`
}

type ChunkResponse struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Text     string  `json:"text"`
	Source   string  `json:"source"`
	StartSec float64 `json:"start_sec"`
	Score    float32 `json:"score"`
}

type SearchResponse struct {
	Query           string          `json:"query"`
	Results         []ChunkResponse `json:"results"`
	Total           int             `json:"total"`
	GeneratedAnswer string          `json:"generated_answer,omitempty"`
	CostoUSD        float64         `json:"costo_usd,omitempty"`
}

type StatsResponse struct {
	IndexName     string `json:"index_name"`
	TotalVectores uint32 `json:"total_vectores"`
	Dimension     int    `json:"dimension"`
	Modelo        string `json:"modelo"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

type OpenAIEmbeddingRequest struct {
	Input      string `json:"input"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions,omitempty"`
}

type OpenAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type OpenAIChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
