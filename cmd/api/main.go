package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/dependencies"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/handlers"
)

func main() {
	log.Println("🚀 Iniciando API de Transcripción...")

	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Inicializar dependencias
	deps, err := dependencies.NewDependencies(cfg)
	if err != nil {
		log.Fatalf("❌ Error inicializando dependencias: %v", err)
	}

	// Inicializar use cases
	appUsecases := NewUsecases(deps, cfg)

	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Configurar rutas
	setupRoutes(r, appUsecases)

	log.Println("✅ Iniciando servidor en el puerto " + cfg.Port)
	// Iniciar servidor
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

// setupRoutes configura todas las rutas de la API
func setupRoutes(r *gin.Engine, usecases Usecases) {
	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Grupo de rutas con prefijo /api/transcribe
	api := r.Group("/api/transcribe")
	{
		api.GET("/health", handlers.HealthCheck(usecases.HealthUseCase))
		api.GET("/stats", handlers.GetStats(usecases.StatsUseCase))
		api.GET("/videos", handlers.GetVideos(usecases.VideoUseCase))
		api.GET("/video/:filename/thumbnail", handlers.ServeThumbnail(usecases.VideoUseCase))
		api.GET("/video/:filename/subtitles", handlers.ServeSubtitles(usecases.VideoUseCase))
		api.GET("/video/:filename", handlers.ServeVideo(usecases.VideoUseCase))
		api.POST("/buscar", handlers.Buscar(usecases.SearchUseCase))
	}
}
