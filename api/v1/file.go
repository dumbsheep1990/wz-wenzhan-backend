package v1

import (
	"net/http"
	"strconv"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// FileAPI 文件上传API处理器
type FileAPI struct {
	fileService service.FileService
}

// NewFileAPI 创建文件上传API处理器
func NewFileAPI(fileService service.FileService) *FileAPI {
	return &FileAPI{
		fileService: fileService,
	}
}

// Upload 上传文件
// @Summary 上传文件
// @Description 上传文件到服务器并返回文件信息
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "上传的文件"
// @Param folder_id formData int false "文件夹ID"
// @Success 201 {object} model.Response{data=model.FileResponse} "上传成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /files/upload [post]
func (api *FileAPI) Upload(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "获取上传文件失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取文件夹ID参数
	var folderID uint
	folderIDStr := c.PostForm("folder_id")
	if folderIDStr != "" {
		folderIDUint, err := strconv.ParseUint(folderIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{
				Code:    http.StatusBadRequest,
				Message: "无效的文件夹ID",
			})
			return
		}
		folderID = uint(folderIDUint)
	}

	// 上传文件
	fileInfo, err := api.fileService.Upload(userID, folderID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "上传文件失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Code:    http.StatusCreated,
		Message: "上传成功",
		Data:    fileInfo,
	})
}

// Delete 删除文件
// @Summary 删除文件
// @Description 删除指定ID的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fileId path int true "文件ID"
// @Success 200 {object} model.Response "删除成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文件不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /files/{fileId} [delete]
func (api *FileAPI) Delete(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	fileIDStr := c.Param("fileId")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的文件ID",
		})
		return
	}

	err = api.fileService.Delete(uint(fileID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "删除文件失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}

// GetURL 获取文件URL
// @Summary 获取文件URL
// @Description 获取指定文件的访问URL
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fileId path int true "文件ID"
// @Success 200 {object} model.Response{data=model.FileURLResponse} "获取成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "文件不存在"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /files/{fileId}/url [get]
func (api *FileAPI) GetURL(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	fileIDStr := c.Param("fileId")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "无效的文件ID",
		})
		return
	}

	url, err := api.fileService.GetURL(uint(fileID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文件URL失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data: model.FileURLResponse{
			URL: url,
		},
	})
}
