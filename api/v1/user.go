package v1

import (
	"net/http"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// UserAPI 用户API处理器
type UserAPI struct {
	userService service.UserService
}

// NewUserAPI 创建用户API处理器
func NewUserAPI(userService service.UserService) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

// Register 用户注册API
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册请求参数"
// @Success 201 {object} model.Response{data=model.UserResponse} "注册成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /users/register [post]
func (api *UserAPI) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	user, err := api.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "注册失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Code:    http.StatusCreated,
		Message: "注册成功",
		Data:    user,
	})
}

// Login 用户登录API
// @Summary 用户登录
// @Description 用户登录并获取授权令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录请求参数"
// @Success 200 {object} model.Response{data=model.LoginResponse} "登录成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "用户名或密码错误"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /users/login [post]
func (api *UserAPI) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	token, user, err := api.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户名或密码错误",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "登录成功",
		Data: model.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

// GetProfile 获取用户信息API
// @Summary 获取用户信息
// @Description 获取当前登录用户的个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Response{data=model.UserResponse} "获取成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /users/profile [get]
func (api *UserAPI) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	user, err := api.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取用户信息失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    user,
	})
}

// UpdateProfile 更新用户信息API
// @Summary 更新用户信息
// @Description 更新当前登录用户的个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UpdateProfileRequest true "更新用户信息请求"
// @Success 200 {object} model.Response{data=model.UserResponse} "更新成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /users/profile [put]
func (api *UserAPI) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户未认证",
		})
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	user, err := api.userService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "更新用户信息失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    user,
	})
}
