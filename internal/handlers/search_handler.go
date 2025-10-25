package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/models"
	"github.com/ngrendenebos/scripts/transcribe-api/internal/usecases"
)

// Buscar retorna un handler para búsqueda
func Buscar(searchUseCase usecases.SearchUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.BusquedaRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Query requerido",
			})
			return
		}

		// Establecer top_k por defecto si no se especifica
		if req.TopK == 0 {
			req.TopK = 10 // Valor por defecto
		}

		// Realizar búsqueda usando el use case
		response, err := searchUseCase.Search(req.Query, req.TopK)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Error en búsqueda",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
