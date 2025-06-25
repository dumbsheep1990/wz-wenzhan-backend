package v1

import (
	"net/http"
	"strconv"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// RecycleAPI 回收站API处理器
type RecycleAPI struct {
	recycleService service.RecycleService
}

// NewRecycleAPI 创建回收站API处理器
func NewRecycleAPI(recycleService service.RecycleService) *RecycleAPI {
	return &RecycleAPI{
		recycleService: recycleService,
	}
}

// List 获取回收站列表
// @Summary 获取回收站列表
// @Description 获取用户回收站中的文档列表，支持分页
// @Tags 回收站管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10"
// @Success 200 {object} model.Response{data=model.RecycleListResponse} "获取成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /recycle [get]
func (api *RecycleAPI) List(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.ListRequest
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

	recycledItems, total, err := api.recycleService.List(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取回收站列表失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data: model.RecycleListResponse{
			Items: recycledItems,
			Total: total,
			Page:  req.Page,
			Size:  req.PageSize,
		},
	})
}

// Restore 恢复回收站项目
// @Summary 恢复回收站项目
// @Description 将回收站中的项目恢复到原位置
// @Tags 回收站管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "回收站项目ID"
// @Success 200 {object} model.Response "恢复成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "项目不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /recycle/{id}/restore [post]
func (api *RecycleAPI) Restore(c *gin.Context) {
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
			Message: "无效的回收站项目ID",
		})
		return
	}

	err = api.recycleService.Restore(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "恢复失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "恢复成功",
	})
}

// DeletePermanently 永久删除回收站项目
// @Summary 永久删除回收站项目
// @Description 从回收站中永久删除项目
// @Tags 回收站管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "回收站项目ID"
// @Success 200 {object} model.Response "删除成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "项目不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /recycle/{id} [delete]
func (api *RecycleAPI) DeletePermanently(c *gin.Context) {
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
			Message: "无效的回收站项目ID",
		})
		return
	}

	err = api.recycleService.DeletePermanently(uint(id), userID)
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

// DeleteBatch 批量永久删除回收站项目
// @Summary 批量永久删除回收站项目
// @Description 从回收站中批量永久删除项目
// @Tags 回收站管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.BatchDeleteRequest true "批量删除请求"
// @Success 200 {object} model.Response "删除成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /recycle/batch [delete]
func (api *RecycleAPI) DeleteBatch(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	err := api.recycleService.DeleteBatch(req.IDs, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "批量删除失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "批量删除成功",
	})
}
