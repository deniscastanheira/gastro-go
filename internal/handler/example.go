package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ExampleHandler demonstra a estrutura de um handler
// Handlers recebem requisições HTTP, parseiam DTOs, chamam Use Cases e retornam respostas
type ExampleHandler struct {
	// Dependências: Use Cases serão injetados aqui
}

// NewExampleHandler cria uma nova instância do handler
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// GetExample demonstra um endpoint GET
func (h *ExampleHandler) GetExample(c echo.Context) error {
	// Parse input (se necessário)
	// Call Use Case
	// Return response
	return c.JSON(http.StatusOK, map[string]string{
		"message": "example endpoint",
	})
}

