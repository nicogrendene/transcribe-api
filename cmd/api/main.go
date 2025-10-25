package main

import (
	"log"
	"net/http"

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
	usecases := NewUsecases(deps, cfg)

	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Configurar rutas
	setupRoutes(r, usecases)

	// Iniciar servidor
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
	log.Println("✅ API lista → http://localhost:" + cfg.Port)
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

	// Rutas
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/health")
	})
	r.GET("/health", handlers.HealthCheck(usecases.HealthUseCase))
	r.GET("/stats", handlers.GetStats(usecases.StatsUseCase))
	r.GET("/videos", handlers.GetVideos(usecases.VideoUseCase))
	r.GET("/video/:filename/thumbnail", handlers.ServeThumbnail(usecases.VideoUseCase))
	r.GET("/video/:filename/subtitles", handlers.ServeSubtitles(usecases.VideoUseCase))
	r.GET("/video/:filename", handlers.ServeVideo(usecases.VideoUseCase))
	r.POST("/buscar", handlers.Buscar(usecases.SearchUseCase))
}
