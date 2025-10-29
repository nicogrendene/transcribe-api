package main

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/config"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/dependencies"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/log"
	"github.com/ngrendenebos/scripts/transcribe-api/cmd/api/middleware"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/handlers"
)

func main() {
	// Inicializar logger
	logger := log.Initialize()
	log.DefaultLogger = logger

	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(context.Background(), "Error cargando configuración", log.Err(err))
	}

	// Inicializar dependencias
	deps, err := dependencies.NewDependencies(cfg)
	if err != nil {
		log.Fatal(context.Background(), "Error inicializando dependencia", log.Err(err))
	}

	// Inicializar use cases
	appUsecases := NewUsecases(deps, cfg)

	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)

	// Deshabilitar logs de Gin
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	gin.DisableConsoleColor()

	// Usar gin.New() en lugar de gin.Default() para evitar middlewares automáticos
	r := gin.New()

	// Configurar middlewares personalizados
	r.Use(middleware.RequestLoggingMiddleware())
	r.Use(middleware.RecoveryWithLogging())

	// Configurar rutas
	setupRoutes(r, appUsecases)

	log.Info(context.Background(), "Iniciando servidor en el puerto", log.Any("port", cfg.Port))
	// Iniciar servidor
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(context.Background(), "Error iniciando servidor", log.Err(err))
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

	// Middleware de compresión gzip
	r.Use(middleware.GzipMiddleware())

	// Middleware de cache para thumbnails
	r.Use(middleware.ThumbnailCacheMiddleware())

	// Middleware de cache para videos
	r.Use(middleware.VideoCacheMiddleware())

	r.GET("/health", handlers.HealthCheck(usecases.HealthUseCase))
	r.GET("/stats", handlers.GetStats(usecases.StatsUseCase))
	r.GET("/videos", handlers.GetVideos(usecases.VideoUseCase))
	r.GET("/video/:id/metadata", handlers.GetVideoMetadata(usecases.VideoUseCase))
	r.GET("/video/:id/info", handlers.GetVideoInfo(usecases.VideoUseCase))
	r.GET("/video/:id/thumbnail", handlers.ServeThumbnail(usecases.VideoUseCase))
	r.GET("/video/:id/subtitles", handlers.ServeSubtitles(usecases.VideoUseCase))
	r.GET("/video/:id/summary", handlers.ServeSummary(usecases.VideoUseCase))
	r.GET("/video/:id", handlers.ServeVideo(usecases.VideoUseCase))
	r.POST("/search", handlers.Search(usecases.SearchUseCase))

}
