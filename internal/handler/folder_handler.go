package handler

import (
	"net/http"
	"strconv"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type FolderHandler struct {
	folderService service.FolderService
}

func NewFolderHandler(folderService service.FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
	}
}

func (h *FolderHandler) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	var req model.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "请求参数错误"))
		return
	}

	folder, err := h.folderService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(folder))
}

func (h *FolderHandler) GetByID(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "无效的文件夹ID"))
		return
	}

	folder, err := h.folderService.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NewErrorResponse(404, "文件夹不存在"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(folder))
}

func (h *FolderHandler) Update(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "无效的文件夹ID"))
		return
	}

	var req model.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "请求参数错误"))
		return
	}

	err = h.folderService.Update(uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(nil))
}

func (h *FolderHandler) Delete(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "无效的文件夹ID"))
		return
	}

	err = h.folderService.Delete(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(nil))
}

func (h *FolderHandler) GetTree(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	tree, err := h.folderService.GetFolderTree(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, "获取文件夹树失败"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(tree))
}

func (h *FolderHandler) GetSubFolders(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	parentIDStr := c.Param("parentId")
	parentID, err := strconv.ParseUint(parentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "无效的父文件夹ID"))
		return
	}

	folders, err := h.folderService.GetSubFolders(uint(parentID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, "获取子文件夹失败"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(folders))
}

func (h *FolderHandler) Move(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(401, "用户未认证"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "无效的文件夹ID"))
		return
	}

	var req model.MoveFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "请求参数错误"))
		return
	}

	err = h.folderService.MoveFolder(uint(id), req.NewParentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(nil))
}
