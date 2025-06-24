package handler

import (
	"net/http"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	// 获取上传的文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "请选择要上传的文件"))
		return
	}

	response, err := h.fileService.UploadFile(userID, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(response))
}

func (h *FileHandler) Delete(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	fileID := c.Param("fileId")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "文件ID不能为空"))
		return
	}

	err := h.fileService.DeleteFile(userID, fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(nil))
}

func (h *FileHandler) GetURL(c *gin.Context) {
	fileID := c.Param("fileId")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "文件ID不能为空"))
		return
	}

	url := h.fileService.GetFileURL(fileID)
	
	response := map[string]string{
		"url": url,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}
