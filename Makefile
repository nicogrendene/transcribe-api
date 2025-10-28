# 🎬 API de Transcripción con Búsqueda Vectorial
# Makefile para facilitar el desarrollo

.PHONY: help build run test clean install deps lint format

# Variables
BINARY_NAME=transcribe-api
BUILD_DIR=bin
MAIN_PATH=cmd/api/main.go
USECASES_PATH=cmd/api/usecases.go

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

## help: Mostrar esta ayuda
help:
	@echo "$(GREEN)🎬 API de Transcripción - Comandos disponibles:$(NC)"
	@echo ""
	@echo "$(YELLOW)Desarrollo:$(NC)"
	@echo "  make run          - Ejecutar la aplicación"
	@echo "  make run-bg       - Ejecutar en background con logs"
	@echo "  make build        - Compilar la aplicación"
	@echo "  make test         - Ejecutar tests"
	@echo "  make lint        - Ejecutar linter"
	@echo "  make format      - Formatear código"
	@echo ""
	@echo "$(YELLOW)Dependencias:$(NC)"
	@echo "  make deps         - Instalar dependencias"
	@echo "  make install     - Instalar en sistema"
	@echo ""
	@echo "$(YELLOW)Utilidades:$(NC)"
	@echo "  make clean       - Limpiar archivos generados"
	@echo "  make setup       - Configurar proyecto inicial"
	@echo ""

## deps: Instalar dependencias
deps:
	@echo "$(GREEN)📦 Instalando dependencias...$(NC)"
	go mod tidy
	go mod download

## build: Compilar la aplicación
build: deps
	@echo "$(GREEN)🔨 Compilando aplicación...$(NC)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH) $(USECASES_PATH)
	@echo "$(GREEN)✅ Compilación exitosa: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## run: Ejecutar la aplicación
run: deps
	@echo "$(GREEN)🚀 Ejecutando API...$(NC)"
	go run $(MAIN_PATH) $(USECASES_PATH)

## run-bg: Ejecutar en background con logs
run-bg: deps
	@echo "$(GREEN)🚀 Ejecutando API en background...$(NC)"
	@echo "$(GREEN)📝 Logs se guardarán en: transcribe-api.log$(NC)"
	@echo "=========================================" >> transcribe-api.log
	@echo "🚀 Nueva sesión iniciada: $$(date)" >> transcribe-api.log
	@echo "=========================================" >> transcribe-api.log
	nohup go run $(MAIN_PATH) $(USECASES_PATH) >> transcribe-api.log 2>&1 &
	@echo "$(GREEN)✅ API ejecutándose en background$(NC)"
	@echo "$(YELLOW)💡 Para ver logs: tail -f transcribe-api.log$(NC)"

## test: Ejecutar tests
test:
	@echo "$(GREEN)🧪 Ejecutando tests...$(NC)"
	go test -v ./...

## lint: Ejecutar linter
lint:
	@echo "$(GREEN)🔍 Ejecutando linter...$(NC)"
	golangci-lint run

## format: Formatear código
format:
	@echo "$(GREEN)✨ Formateando código...$(NC)"
	go fmt ./...
	goimports -w .

## install: Instalar en sistema
install: build
	@echo "$(GREEN)📦 Instalando en sistema...$(NC)"
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)✅ Instalado en /usr/local/bin/$(BINARY_NAME)$(NC)"

## clean: Limpiar archivos generados
clean:
	@echo "$(GREEN)🧹 Limpiando archivos...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f transcribe-api.log
	go clean
	@echo "$(GREEN)✅ Limpieza completada$(NC)"

## setup: Configurar proyecto inicial
setup:
	@echo "$(GREEN)⚙️  Configurando proyecto...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)📝 Creando archivo .env...$(NC)"; \
		cp env.example .env; \
		echo "$(GREEN)✅ Archivo .env creado. Edita con tus API keys.$(NC)"; \
	else \
		echo "$(GREEN)✅ Archivo .env ya existe$(NC)"; \
	fi
	@echo "$(GREEN)📦 Instalando dependencias...$(NC)"
	$(MAKE) deps
	@echo "$(GREEN)🎉 Configuración completada!$(NC)"
	@echo "$(YELLOW)💡 Próximos pasos:$(NC)"
	@echo "  1. Edita .env con tus API keys"
	@echo "  2. Ejecuta: make run"

## dev: Modo desarrollo con auto-reload
dev:
	@echo "$(GREEN)🔄 Modo desarrollo (requiere air)...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)❌ Air no está instalado. Instala con: go install github.com/cosmtrek/air@latest$(NC)"; \
	fi

## docker: Construir imagen Docker
docker:
	@echo "$(GREEN)🐳 Construyendo imagen Docker...$(NC)"
	docker build -t $(BINARY_NAME) .

## docker-run: Ejecutar con Docker
docker-run: docker
	@echo "$(GREEN)🚀 Ejecutando con Docker...$(NC)"
	docker run -p 8000:8000 --env-file .env $(BINARY_NAME)

# Comando por defecto
.DEFAULT_GOAL := help
