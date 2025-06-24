package main

import (
	"log"
	"wz-wenzhan-backend/internal/config"
	"wz-wenzhan-backend/internal/handler"
	"wz-wenzhan-backend/internal/middleware"
	"wz-wenzhan-backend/internal/repository"
	"wz-wenzhan-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化数据库
	db := config.InitDB(cfg)

	// 初始化Redis
	_ = config.InitRedis(cfg)

	// 初始化日志
	logger := config.InitLogger(cfg)

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	documentRepo := repository.NewDocumentRepository(db)
	folderRepo := repository.NewFolderRepository(db)
	workspaceRepo := repository.NewWorkspaceRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	recycleRepo := repository.NewRecycleRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, logger)
	documentService := service.NewDocumentService(documentRepo, logger)
	folderService := service.NewFolderService(folderRepo, documentRepo, activityRepo, logger)
	fileService := service.NewFileService("./uploads", "http://localhost:8080", logger)
	searchService := service.NewSearchService(documentRepo, folderRepo, logger)
	workspaceService := service.NewWorkspaceService(workspaceRepo, logger)
	activityService := service.NewActivityService(activityRepo, logger)
	recycleService := service.NewRecycleService(recycleRepo, logger)

	// 初始化处理器层
	userHandler := handler.NewUserHandler(userService)
	documentHandler := handler.NewDocumentHandler(documentService)
	folderHandler := handler.NewFolderHandler(folderService)
	fileHandler := handler.NewFileHandler(fileService)
	searchHandler := handler.NewSearchHandler(searchService)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceService)
	activityHandler := handler.NewActivityHandler(activityService)
	recycleHandler := handler.NewRecycleHandler(recycleService)

	// 初始化Gin引擎
	r := gin.Default()

	// 设置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 设置中间件
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))

	// 静态文件服务（用于文件上传）
	r.Static("/uploads", "./uploads")

	// 注册路由
	setupRoutes(r, userHandler, documentHandler, folderHandler, fileHandler, 
		searchHandler, workspaceHandler, activityHandler, recycleHandler)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine,
	userHandler *handler.UserHandler,
	documentHandler *handler.DocumentHandler,
	folderHandler *handler.FolderHandler,
	fileHandler *handler.FileHandler,
	searchHandler *handler.SearchHandler,
	workspaceHandler *handler.WorkspaceHandler,
	activityHandler *handler.ActivityHandler,
	recycleHandler *handler.RecycleHandler) {

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

	// 足迹记录相关路由
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
