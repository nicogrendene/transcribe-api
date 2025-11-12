package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
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
	videos    map[string]models.Video // Mapa de videos por ID
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

	log.Info(context.Background(), "Conectado a Pinecone", log.Any("index", index), log.Any("stats", stats))

	// Cargar videos desde el archivo JSON
	videos, err := loadVideos()
	if err != nil {
		log.Error(context.Background(), "Error cargando videos.json", log.Err(err))
		videos = make(map[string]models.Video) // Continuar con mapa vacío
	}

	return &PineconeService{
		Client:    client,
		Index:     index,
		IndexName: indexName,
		videos:    videos,
	}, nil
}

// loadVideos carga los videos desde el archivo videos.json
func loadVideos() (map[string]models.Video, error) {
	data, err := os.ReadFile("videos.json")
	if err != nil {
		return nil, fmt.Errorf("error leyendo videos.json: %v", err)
	}

	var videosData models.VideosData
	if err := json.Unmarshal(data, &videosData); err != nil {
		return nil, fmt.Errorf("error parseando videos.json: %v", err)
	}

	// Crear mapa de videos por ID
	videoMap := make(map[string]models.Video)
	for _, video := range videosData.Videos {
		videoMap[video.ID] = video
	}

	return videoMap, nil
}

// Search realiza una búsqueda vectorial en Pinecone
func (s *PineconeService) Search(ctx context.Context, embedding []float32, topK int) ([]models.ChunkResponse, error) {
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

func (s *PineconeService) parseResults(matches []*pinecone.ScoredVector) []models.ChunkResponse {
	var res []models.ChunkResponse

	for _, match := range matches {
		chunk := s.extractMetadata(match)
		chunk.Score = match.Score
		res = append(res, chunk)
	}

	return res
}

func (s *PineconeService) extractMetadata(match *pinecone.ScoredVector) models.ChunkResponse {
	var chunk models.ChunkResponse

	// Extraer el ID del vector
	chunk.ID = match.Vector.Id

	// Variable temporal para guardar el source_file
	var sourceFile string

	if match.Vector.Metadata != nil && match.Vector.Metadata.Fields != nil {
		for key, val := range match.Vector.Metadata.Fields {
			if kind := val.GetKind(); kind != nil {
				if key == "source_file" {
					if v, ok := kind.(*structpb.Value_StringValue); ok {
						sourceFile = utils.CleanPointerFormat(v.StringValue)
					}
				}
				s.setMetadataField(&chunk, key, kind)
			}
		}
	}

	// Asignar el video con el source_file
	chunk.Video = sourceFile

	// Buscar el video en el mapa y asignar el source y url
	if video, exists := s.videos[sourceFile]; exists {
		chunk.Source = video.Source
		chunk.URL = video.URL
	}

	return chunk
}

func (s *PineconeService) setMetadataField(chunk *models.ChunkResponse, key string, kind interface{}) {
	switch key {
	case "title":
		if v, ok := kind.(*structpb.Value_StringValue); ok {
			chunk.Title = utils.CleanPointerFormat(v.StringValue)
		}
	case "text":
		if v, ok := kind.(*structpb.Value_StringValue); ok {
			chunk.Text = utils.CleanPointerFormat(v.StringValue)
		}
	case "start_sec":
		chunk.StartSec = s.parseFloatFromMetadata(kind)
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
func (s *PineconeService) GetStats(ctx context.Context) (*models.StatsResponse, error) {
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
