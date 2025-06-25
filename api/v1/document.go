package v1

import (
	"net/http"
	"strconv"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// DocumentAPI 文档API处理器
type DocumentAPI struct {
	documentService service.DocumentService
}

// NewDocumentAPI 创建文档API处理器
func NewDocumentAPI(documentService service.DocumentService) *DocumentAPI {
	return &DocumentAPI{
		documentService: documentService,
	}
}

// List 获取文档列表
// @Summary 获取文档列表
// @Description 获取用户的文档列表，支持分页和过滤
// @Tags 文档管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10"
// @Param folder_id query int false "文件夹ID"
// @Param type query string false "文档类型：word, excel, mindmap, note, ai_draft, imported"
// @Param status query int false "文档状态：1-草稿，2-已发布，3-已归档"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} model.Response{data=model.DocumentListResponse} "获取成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /documents [get]
func (api *DocumentAPI) List(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.DocumentListRequest
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

	documents, total, err := api.documentService.List(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文档列表失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data: model.DocumentListResponse{
			Items: documents,
			Total: total,
			Page:  req.Page,
			Size:  req.PageSize,
		},
	})
}

// Create 创建文档
// @Summary 创建文档
// @Description 创建新文档
// @Tags 文档管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateDocumentRequest true "创建文档请求"
// @Success 201 {object} model.Response{data=model.DocumentResponse} "创建成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /documents [post]
func (api *DocumentAPI) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	document, err := api.documentService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "创建文档失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    document,
	})
}

// GetByID 获取文档详情
// @Summary 获取文档详情
// @Description 获取指定ID的文档详细信息
// @Tags 文档管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文档ID"
// @Success 200 {object} model.Response{data=model.DocumentResponse} "获取成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文档不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /documents/{id} [get]
func (api *DocumentAPI) GetByID(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的文档ID",
		})
		return
	}

	document, err := api.documentService.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    http.StatusNotFound,
			Message: "文档不存在",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    document,
	})
}

// Update 更新文档
// @Summary 更新文档
// @Description 更新指定ID的文档
// @Tags 文档管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文档ID"
// @Param request body model.UpdateDocumentRequest true "更新文档请求"
// @Success 200 {object} model.Response "更新成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文档不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /documents/{id} [put]
func (api *DocumentAPI) Update(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的文档ID",
		})
		return
	}

	var req model.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	err = api.documentService.Update(uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "更新失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
	})
}

// Delete 删除文档
// @Summary 删除文档
// @Description 删除指定ID的文档（移入回收站）
// @Tags 文档管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文档ID"
// @Success 200 {object} model.Response "删除成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文档不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /documents/{id} [delete]
func (api *DocumentAPI) Delete(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的文档ID",
		})
		return
	}

	err = api.documentService.Delete(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "删除失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}

// Share 分享文档
// @Summary 分享文档
// @Description 创建文档分享链接
// @Tags 文档管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文档ID"
// @Param request body model.ShareDocumentRequest true "分享文档请求"
// @Success 200 {object} model.Response{data=model.ShareResponse} "分享成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文档不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /documents/{id}/share [post]
func (api *DocumentAPI) Share(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的文档ID",
		})
		return
	}

	var req model.ShareDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	token, err := api.documentService.Share(uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "分享失败",
			Error:   err.Error(),
		})
		return
	}

	// 使用域名和端口构建完整的分享URL
	shareURL := "/api/v1/share/" + token

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "分享成功",
		Data: model.ShareResponse{
			Token:    token,
			ShareURL: shareURL,
		},
	})
}

// Copy 复制文档
// @Summary 复制文档
// @Description 复制指定ID的文档
// @Tags 文档管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文档ID"
// @Success 201 {object} model.Response{data=model.DocumentResponse} "复制成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文档不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /documents/{id}/copy [post]
func (api *DocumentAPI) Copy(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的文档ID",
		})
		return
	}

	document, err := api.documentService.Copy(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "复制失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Code:    http.StatusCreated,
		Message: "复制成功",
		Data:    document,
	})
}
