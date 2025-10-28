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

	r.GET("/health", handlers.HealthCheck(usecases.HealthUseCase))
	r.GET("/stats", handlers.GetStats(usecases.StatsUseCase))
	r.GET("/videos", handlers.GetVideos(usecases.VideoUseCase))
	r.GET("/video/:id/thumbnail", handlers.ServeThumbnail(usecases.VideoUseCase))
	r.GET("/video/:id/subtitles", handlers.ServeSubtitles(usecases.VideoUseCase))
	r.GET("/video/:id/summary", handlers.ServeSummary(usecases.VideoUseCase))
	r.GET("/video/:id", handlers.ServeVideo(usecases.VideoUseCase))
	r.POST("/search", handlers.Search(usecases.SearchUseCase))

}
