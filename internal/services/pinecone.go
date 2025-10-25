package services

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/pkg/utils"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
)

// PineconeService maneja las interacciones con Pinecone
type PineconeService struct {
	Client    *pinecone.Client
	Index     *pinecone.IndexConnection
	IndexName string
}

// NewPineconeService crea una nueva instancia del servicio Pinecone
func NewPineconeService(apiKey, indexName string) (*PineconeService, error) {
	ctx := context.Background()

	client, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("error inicializando Pinecone: %v", err)
	}

	indexList, err := client.ListIndexes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listando índices: %v", err)
	}

	var indexHost string
	for _, idx := range indexList {
		if idx.Name == indexName {
			indexHost = idx.Host
			break
		}
	}

	if indexHost == "" {
		return nil, fmt.Errorf("índice '%s' no encontrado", indexName)
	}

	index, err := client.Index(pinecone.NewIndexConnParams{
		Host: indexHost,
	})
	if err != nil {
		return nil, fmt.Errorf("error conectando al índice: %v", err)
	}

	stats, err := index.DescribeIndexStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo stats: %v", err)
	}

	log.Printf("✅ Índice '%s' conectado (%d vectores)", indexName, stats.TotalVectorCount)

	return &PineconeService{
		Client:    client,
		Index:     index,
		IndexName: indexName,
	}, nil
}

// Search realiza una búsqueda vectorial en Pinecone
func (s *PineconeService) Search(embedding []float32, topK int) ([]models.ChunkResultado, error) {
	ctx := context.Background()
	queryReq := &pinecone.QueryByVectorValuesRequest{
		Vector:          embedding,
		TopK:            uint32(topK),
		IncludeMetadata: true,
	}

	results, err := s.Index.QueryByVectorValues(ctx, queryReq)
	if err != nil {
		return nil, fmt.Errorf("error en búsqueda: %v", err)
	}

	return s.parseResults(results.Matches), nil
}

// parseResults convierte los resultados de Pinecone en ChunkResultado
func (s *PineconeService) parseResults(matches []*pinecone.ScoredVector) []models.ChunkResultado {
	var resultados []models.ChunkResultado

	for _, match := range matches {
		chunk := s.extractMetadata(match)
		chunk.Score = match.Score
		resultados = append(resultados, chunk)
	}

	return resultados
}

// extractMetadata extrae los metadatos de un vector de Pinecone
func (s *PineconeService) extractMetadata(match *pinecone.ScoredVector) models.ChunkResultado {
	var chunk models.ChunkResultado

	if match.Vector.Metadata != nil && match.Vector.Metadata.Fields != nil {
		for key, val := range match.Vector.Metadata.Fields {
			if kind := val.GetKind(); kind != nil {
				s.setMetadataField(&chunk, key, kind)
			}
		}
	}

	return chunk
}

// setMetadataField establece un campo específico de los metadatos
func (s *PineconeService) setMetadataField(chunk *models.ChunkResultado, key string, kind interface{}) {
	switch key {
	case "chunk_id":
		if v, ok := kind.(*structpb.Value_StringValue); ok {
			chunk.ChunkID = utils.CleanPointerFormat(v.StringValue)
		}
	case "title":
		if v, ok := kind.(*structpb.Value_StringValue); ok {
			chunk.Title = utils.CleanPointerFormat(v.StringValue)
		}
	case "text":
		if v, ok := kind.(*structpb.Value_StringValue); ok {
			chunk.Text = utils.CleanPointerFormat(v.StringValue)
		}
	case "source_file":
		if v, ok := kind.(*structpb.Value_StringValue); ok {
			chunk.SourceFile = utils.CleanPointerFormat(v.StringValue)
		}
	case "start_sec":
		chunk.StartSec = s.parseFloatFromMetadata(kind)
	case "end_sec":
		chunk.EndSec = s.parseFloatFromMetadata(kind)
	}
}

// parseFloatFromMetadata convierte un valor de metadatos a float64
func (s *PineconeService) parseFloatFromMetadata(kind interface{}) float64 {
	if v, ok := kind.(*structpb.Value_StringValue); ok {
		cleanVal := utils.CleanPointerFormat(v.StringValue)
		if f, err := strconv.ParseFloat(cleanVal, 64); err == nil {
			return f
		}
	} else if v, ok := kind.(*structpb.Value_NumberValue); ok {
		return v.NumberValue
	}
	return 0
}

// GetStats obtiene las estadísticas del índice
func (s *PineconeService) GetStats() (*models.StatsResponse, error) {
	ctx := context.Background()
	stats, err := s.Index.DescribeIndexStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo stats: %v", err)
	}

	return &models.StatsResponse{
		IndexName:     s.IndexName,
		TotalVectores: stats.TotalVectorCount,
		Dimension:     512,                      // Usar constante del config
		Modelo:        "text-embedding-3-small", // Usar constante del config
	}, nil
}

// cleanPointerFormat limpia el formato de puntero "&{valor}" que puede venir de Pinecone
func cleanPointerFormat(s string) string {
	if len(s) > 3 && s[:2] == "&{" && s[len(s)-1] == '}' {
		return s[2 : len(s)-1]
	}
	return s
}
