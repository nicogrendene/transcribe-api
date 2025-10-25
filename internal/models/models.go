package models

// Request/Response models

// BusquedaRequest representa la petición de búsqueda
type BusquedaRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

// ChunkResultado representa un fragmento de video encontrado
type ChunkResultado struct {
	ChunkID    string  `json:"chunk_id"`
	Title      string  `json:"title"`
	Text       string  `json:"text"`
	StartSec   float64 `json:"start_sec"`
	EndSec     float64 `json:"end_sec"`
	SourceFile string  `json:"source_file"`
	Score      float32 `json:"score"`
}

// BusquedaResponse representa la respuesta de búsqueda
type BusquedaResponse struct {
	Query      string           `json:"query"`
	Resultados []ChunkResultado `json:"resultados"`
	Total      int              `json:"total"`
	CostoUSD   float64          `json:"costo_usd,omitempty"`
}

// StatsResponse representa las estadísticas del índice
type StatsResponse struct {
	IndexName     string `json:"index_name"`
	TotalVectores uint32 `json:"total_vectores"`
	Dimension     int    `json:"dimension"`
	Modelo        string `json:"modelo"`
}

// HealthResponse representa el estado de salud de la API
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// OpenAI API models

// OpenAIEmbeddingRequest representa la petición a OpenAI Embeddings API
type OpenAIEmbeddingRequest struct {
	Input      string `json:"input"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions,omitempty"`
}

// OpenAIEmbeddingResponse representa la respuesta de OpenAI Embeddings API
type OpenAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}
