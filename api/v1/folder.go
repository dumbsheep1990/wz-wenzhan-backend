package v1

import (
	"net/http"
	"strconv"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// FolderAPI 文件夹API处理器
type FolderAPI struct {
	folderService service.FolderService
}

// NewFolderAPI 创建文件夹API处理器
func NewFolderAPI(folderService service.FolderService) *FolderAPI {
	return &FolderAPI{
		folderService: folderService,
	}
}

// GetTree 获取文件夹树
// @Summary 获取文件夹树
// @Description 获取用户的文件夹层级结构
// @Tags 文件夹管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Response{data=[]model.FolderTreeResponse} "获取成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /folders/tree [get]
func (api *FolderAPI) GetTree(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	folders, err := api.folderService.GetTree(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文件夹树失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    folders,
	})
}

// Create 创建文件夹
// @Summary 创建文件夹
// @Description 创建新的文件夹
// @Tags 文件夹管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateFolderRequest true "创建文件夹请求"
// @Success 201 {object} model.Response{data=model.FolderResponse} "创建成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /folders [post]
func (api *FolderAPI) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	folder, err := api.folderService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "创建文件夹失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    folder,
	})
}

// GetByID 获取文件夹详情
// @Summary 获取文件夹详情
// @Description 获取指定ID的文件夹详细信息
// @Tags 文件夹管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件夹ID"
// @Success 200 {object} model.Response{data=model.FolderResponse} "获取成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文件夹不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /folders/{id} [get]
func (api *FolderAPI) GetByID(c *gin.Context) {
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
			Message: "无效的文件夹ID",
		})
		return
	}

	folder, err := api.folderService.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    http.StatusNotFound,
			Message: "文件夹不存在",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    folder,
	})
}

// Update 更新文件夹
// @Summary 更新文件夹
// @Description 更新指定ID的文件夹
// @Tags 文件夹管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件夹ID"
// @Param request body model.UpdateFolderRequest true "更新文件夹请求"
// @Success 200 {object} model.Response "更新成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文件夹不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /folders/{id} [put]
func (api *FolderAPI) Update(c *gin.Context) {
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
			Message: "无效的文件夹ID",
		})
		return
	}

	var req model.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	err = api.folderService.Update(uint(id), userID, &req)
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

// Delete 删除文件夹
// @Summary 删除文件夹
// @Description 删除指定ID的文件夹
// @Tags 文件夹管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件夹ID"
// @Success 200 {object} model.Response "删除成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文件夹不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /folders/{id} [delete]
func (api *FolderAPI) Delete(c *gin.Context) {
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
			Message: "无效的文件夹ID",
		})
		return
	}

	err = api.folderService.Delete(uint(id), userID)
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

// Move 移动文件夹
// @Summary 移动文件夹
// @Description 移动文件夹到新的父文件夹
// @Tags 文件夹管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件夹ID"
// @Param request body model.MoveFolderRequest true "移动文件夹请求"
// @Success 200 {object} model.Response "移动成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文件夹不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /folders/{id}/move [post]
func (api *FolderAPI) Move(c *gin.Context) {
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
			Message: "无效的文件夹ID",
		})
		return
	}

	var req model.MoveFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	err = api.folderService.Move(uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "移动失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "移动成功",
	})
}

// GetSubFolders 获取子文件夹
// @Summary 获取子文件夹
// @Description 获取指定父文件夹下的所有子文件夹
// @Tags 文件夹管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param parentId path int true "父文件夹ID，0表示根目录"
// @Success 200 {object} model.Response{data=[]model.FolderResponse} "获取成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /folders/{parentId}/subfolders [get]
func (api *FolderAPI) GetSubFolders(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	parentIDStr := c.Param("parentId")
	parentID, err := strconv.ParseUint(parentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的父文件夹ID",
		})
		return
	}

	folders, err := api.folderService.GetSubFolders(userID, uint(parentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取子文件夹失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    folders,
	})
}
