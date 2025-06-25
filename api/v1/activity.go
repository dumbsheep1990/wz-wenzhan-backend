package v1

import (
	"net/http"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ActivityAPI 活动记录API处理器
type ActivityAPI struct {
	activityService service.ActivityService
}

// NewActivityAPI 创建活动记录API处理器
func NewActivityAPI(activityService service.ActivityService) *ActivityAPI {
	return &ActivityAPI{
		activityService: activityService,
	}
}

// List 获取活动记录列表
// @Summary 获取活动记录列表
// @Description 获取用户的活动记录列表，支持分页
// @Tags 活动记录管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10"
// @Success 200 {object} model.Response{data=model.ActivityListResponse} "获取成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /activities [get]
func (api *ActivityAPI) List(c *gin.Context) {
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

	activities, total, err := api.activityService.List(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取活动记录失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data: model.ActivityListResponse{
			Items: activities,
			Total: total,
			Page:  req.Page,
			Size:  req.PageSize,
		},
	})
}

// Create 创建活动记录
// @Summary 创建活动记录
// @Description 创建用户活动记录(通常由系统自动创建)
// @Tags 活动记录管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateActivityRequest true "创建活动记录请求"
// @Success 201 {object} model.Response{data=model.ActivityResponse} "创建成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /activities [post]
func (api *ActivityAPI) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	activity, err := api.activityService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "创建活动记录失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    activity,
	})
}
