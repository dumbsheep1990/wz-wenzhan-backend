package handler

import (
	"net/http"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService service.SearchService
}

func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) SearchDocuments(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	var req model.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "请求参数错误"))
		return
	}

	result, err := h.searchService.SearchDocuments(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, "搜索失败"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *SearchHandler) SearchAll(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	var req model.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "请求参数错误"))
		return
	}

	result, err := h.searchService.SearchAll(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, "搜索失败"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}
