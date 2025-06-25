package v1

import (
	"net/http"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// WorkspaceAPI 工作台API处理器
type WorkspaceAPI struct {
	workspaceService service.WorkspaceService
}

// NewWorkspaceAPI 创建工作台API处理器
func NewWorkspaceAPI(workspaceService service.WorkspaceService) *WorkspaceAPI {
	return &WorkspaceAPI{
		workspaceService: workspaceService,
	}
}

// GetDashboard 获取工作台首页信息
// @Summary 获取工作台首页信息
// @Description 获取用户的工作台首页信息，包括最近访问、重要文档等
// @Tags 工作台管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Response{data=model.DashboardResponse} "获取成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /workspace/dashboard [get]
func (api *WorkspaceAPI) GetDashboard(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	// 获取最近访问的文档
	recentDocuments, err := api.workspaceService.GetRecentDocuments(userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取最近访问文档失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取重要文档
	importantDocuments, err := api.workspaceService.GetImportantDocuments(userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取重要文档失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取活动记录
	activities, err := api.workspaceService.GetRecentActivities(userID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取活动记录失败",
			Error:   err.Error(),
		})
		return
	}

	// 构建工作台响应
	dashboard := model.DashboardResponse{
		RecentDocuments:    recentDocuments,
		ImportantDocuments: importantDocuments,
		Activities:         activities,
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    dashboard,
	})
}

// GetStats 获取工作台统计信息
// @Summary 获取工作台统计信息
// @Description 获取用户的工作台统计信息，包括文档数量、存储使用量等
// @Tags 工作台管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Response{data=model.StatsResponse} "获取成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /workspace/stats [get]
func (api *WorkspaceAPI) GetStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	// 获取文档统计
	documentCount, err := api.workspaceService.GetDocumentCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文档数量统计失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取文件夹统计
	folderCount, err := api.workspaceService.GetFolderCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文件夹数量统计失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取存储使用量
	storageUsed, err := api.workspaceService.GetStorageUsed(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取存储使用量失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取类型分布
	typeDistribution, err := api.workspaceService.GetTypeDistribution(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文档类型分布失败",
			Error:   err.Error(),
		})
		return
	}
	
	// 构建统计数据
	stats := model.StatsResponse{
		DocumentCount:    documentCount,
		FolderCount:      folderCount,
		StorageUsed:      storageUsed,
		StorageLimit:     10 * 1024 * 1024 * 1024, // 10GB默认存储空间
		TypeDistribution: typeDistribution,
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    stats,
	})
}
