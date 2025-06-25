package service

import (
	"mime/multipart"
	"wz-wenzhan-backend/internal/model"
)

// UserService 用户服务接口
type UserService interface {
	Register(req *model.RegisterRequest) (*model.UserResponse, error)
	Login(req *model.LoginRequest) (string, *model.UserResponse, error)
	GetByID(userID uint) (*model.UserResponse, error)
	UpdateProfile(userID uint, req *model.UpdateProfileRequest) (*model.UserResponse, error)
}

// DocumentService 文档服务接口
type DocumentService interface {
	Create(userID uint, req *model.CreateDocumentRequest) (*model.DocumentResponse, error)
	GetByID(id uint, userID uint) (*model.DocumentResponse, error)
	Update(id uint, userID uint, req *model.UpdateDocumentRequest) error
	Delete(id uint, userID uint) error
	List(userID uint, req *model.DocumentListRequest) ([]*model.DocumentResponse, int64, error)
	Share(id uint, userID uint, req *model.ShareDocumentRequest) (string, error)
	Copy(id uint, userID uint) (*model.DocumentResponse, error)
}

// FolderService 文件夹服务接口
type FolderService interface {
	Create(userID uint, req *model.CreateFolderRequest) (*model.FolderResponse, error)
	GetByID(id uint, userID uint) (*model.FolderResponse, error)
	Update(id uint, userID uint, req *model.UpdateFolderRequest) error
	Delete(id uint, userID uint) error
	GetTree(userID uint) ([]*model.FolderTreeResponse, error)
	Move(id uint, userID uint, req *model.MoveFolderRequest) error
	GetSubFolders(userID uint, parentID uint) ([]*model.FolderResponse, error)
}

// FileService 文件服务接口
type FileService interface {
	Upload(userID uint, folderID uint, file *multipart.FileHeader) (*model.FileResponse, error)
	Delete(id uint, userID uint) error
	GetURL(id uint, userID uint) (string, error)
}

// SearchService 搜索服务接口
type SearchService interface {
	SearchDocuments(userID uint, req *model.SearchDocumentsRequest) ([]*model.DocumentResponse, int64, error)
	SearchAll(userID uint, req *model.SearchAllRequest) ([]*model.SearchResultItem, int64, error)
}

// WorkspaceService 工作台服务接口
type WorkspaceService interface {
	// 工作台首页接口
	GetRecentDocuments(userID uint, limit int) ([]*model.DocumentResponse, error)
	GetImportantDocuments(userID uint, limit int) ([]*model.DocumentResponse, error)
	GetRecentActivities(userID uint, limit int) ([]*model.ActivityResponse, error)
	
	// 统计信息接口
	GetDocumentCount(userID uint) (int64, error)
	GetFolderCount(userID uint) (int64, error)
	GetStorageUsed(userID uint) (int64, error)
	GetTypeDistribution(userID uint) (map[string]int64, error)
}

// ActivityService 活动记录服务接口
type ActivityService interface {
	Create(userID uint, req *model.CreateActivityRequest) (*model.ActivityResponse, error)
	List(userID uint, req *model.ListRequest) ([]*model.ActivityResponse, int64, error)
}

// RecycleService 回收站服务接口
type RecycleService interface {
	List(userID uint, req *model.ListRequest) ([]*model.RecycleItemResponse, int64, error)
	Restore(id uint, userID uint) error
	DeletePermanently(id uint, userID uint) error
	DeleteBatch(ids []uint, userID uint) error
}
