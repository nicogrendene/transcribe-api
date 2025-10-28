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
	Query    string          `json:"query"`
	Results  []ChunkResponse `json:"results"`
	Total    int             `json:"total"`
	CostoUSD float64         `json:"costo_usd,omitempty"`
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
