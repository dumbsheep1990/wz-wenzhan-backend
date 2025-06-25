package v1

import (
	"wz-wenzhan-backend/internal/handler"
	"wz-wenzhan-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置API路由
func SetupRoutes(r *gin.Engine, 
	userHandler *handler.UserHandler,
	documentHandler *handler.DocumentHandler,
	folderHandler *handler.FolderHandler,
	fileHandler *handler.FileHandler,
	searchHandler *handler.SearchHandler,
	workspaceHandler *handler.WorkspaceHandler,
	activityHandler *handler.ActivityHandler,
	recycleHandler *handler.RecycleHandler) {

	// API v1路由组
	api := r.Group("/api/v1")

	// 用户相关路由
	users := api.Group("/users")
	{
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.Login)
		users.GET("/profile", middleware.AuthRequired(), userHandler.GetProfile)
		users.PUT("/profile", middleware.AuthRequired(), userHandler.UpdateProfile)
	}

	// 工作台相关路由
	workspace := api.Group("/workspace")
	workspace.Use(middleware.AuthRequired())
	{
		workspace.GET("/dashboard", workspaceHandler.GetDashboard)
		workspace.GET("/stats", workspaceHandler.GetStats)
	}

	// 文档相关路由
	documents := api.Group("/documents")
	documents.Use(middleware.AuthRequired())
	{
		documents.GET("", documentHandler.List)
		documents.POST("", documentHandler.Create)
		documents.GET("/:id", documentHandler.GetByID)
		documents.PUT("/:id", documentHandler.Update)
		documents.DELETE("/:id", documentHandler.Delete)
		documents.POST("/:id/share", documentHandler.Share)
		documents.POST("/:id/copy", documentHandler.Copy)
	}

	// 文件夹相关路由
	folders := api.Group("/folders")
	folders.Use(middleware.AuthRequired())
	{
		folders.GET("/tree", folderHandler.GetTree)
		folders.POST("", folderHandler.Create)
		folders.GET("/:id", folderHandler.GetByID)
		folders.PUT("/:id", folderHandler.Update)
		folders.DELETE("/:id", folderHandler.Delete)
		folders.POST("/:id/move", folderHandler.Move)
		folders.GET("/:parentId/subfolders", folderHandler.GetSubFolders)
	}

	// 文件上传相关路由
	files := api.Group("/files")
	files.Use(middleware.AuthRequired())
	{
		files.POST("/upload", fileHandler.Upload)
		files.DELETE("/:fileId", fileHandler.Delete)
		files.GET("/:fileId/url", fileHandler.GetURL)
	}

	// 搜索相关路由
	search := api.Group("/search")
	search.Use(middleware.AuthRequired())
	{
		search.GET("/documents", searchHandler.SearchDocuments)
		search.GET("/all", searchHandler.SearchAll)
	}

	// 活动记录相关路由
	activities := api.Group("/activities")
	activities.Use(middleware.AuthRequired())
	{
		activities.GET("", activityHandler.List)
		activities.POST("", activityHandler.Create)
	}

	// 回收站相关路由
	recycle := api.Group("/recycle")
	recycle.Use(middleware.AuthRequired())
	{
		recycle.GET("", recycleHandler.List)
		recycle.POST("/:id/restore", recycleHandler.Restore)
		recycle.DELETE("/:id", recycleHandler.DeletePermanently)
		recycle.DELETE("/batch", recycleHandler.DeleteBatch)
	}
}
