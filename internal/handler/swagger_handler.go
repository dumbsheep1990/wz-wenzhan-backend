package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SwaggerHandler handles API documentation routes
type SwaggerHandler struct {
	BasePath string
}

// NewSwaggerHandler creates a new swagger handler
func NewSwaggerHandler(basePath string) *SwaggerHandler {
	return &SwaggerHandler{
		BasePath: basePath,
	}
}

// ServeSwaggerUI serves the swagger UI
func (h *SwaggerHandler) ServeSwaggerUI(c *gin.Context) {
	c.HTML(http.StatusOK, "swagger.html", gin.H{
		"title": "API 文档 - 万知文站",
	})
}

// GetSwaggerJSON returns swagger.json file
func (h *SwaggerHandler) GetSwaggerJSON(c *gin.Context) {
	swaggerPath := filepath.Join(h.BasePath, "api/swagger.json")
	c.File(swaggerPath)
}

// GetSwaggerYAML returns swagger.yaml file
func (h *SwaggerHandler) GetSwaggerYAML(c *gin.Context) {
	swaggerPath := filepath.Join(h.BasePath, "api/swagger.yaml")
	c.File(swaggerPath)
}
