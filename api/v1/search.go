package v1

import (
	"net/http"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SearchAPI 搜索API处理器
type SearchAPI struct {
	searchService service.SearchService
}

// NewSearchAPI 创建搜索API处理器
func NewSearchAPI(searchService service.SearchService) *SearchAPI {
	return &SearchAPI{
		searchService: searchService,
	}
}

// SearchDocuments 搜索文档
// @Summary 搜索文档
// @Description 根据关键词搜索文档
// @Tags 搜索管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string true "搜索关键词"
// @Param folder_id query int false "在指定文件夹中搜索，不指定则搜索全部"
// @Param type query string false "文档类型过滤"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10"
// @Success 200 {object} model.Response{data=model.SearchDocumentsResponse} "搜索成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /search/documents [get]
func (api *SearchAPI) SearchDocuments(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.SearchDocumentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	if req.Keyword == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "搜索关键词不能为空",
		})
		return
	}

	documents, total, err := api.searchService.SearchDocuments(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "搜索文档失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "搜索成功",
		Data: model.SearchDocumentsResponse{
			Items: documents,
			Total: total,
			Page:  req.Page,
			Size:  req.PageSize,
		},
	})
}

// SearchAll 全局搜索
// @Summary 全局搜索
// @Description 在全站范围内搜索文档、文件夹等
// @Tags 搜索管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10"
// @Success 200 {object} model.Response{data=model.SearchAllResponse} "搜索成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /search/all [get]
func (api *SearchAPI) SearchAll(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.SearchAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	if req.Keyword == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "搜索关键词不能为空",
		})
		return
	}

	results, total, err := api.searchService.SearchAll(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "全局搜索失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "搜索成功",
		Data: model.SearchAllResponse{
			Items: results,
			Total: total,
			Page:  req.Page,
			Size:  req.PageSize,
		},
	})
}
