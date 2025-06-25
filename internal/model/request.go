package model

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ListRequest 通用列表请求
type ListRequest struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	Title    string       `json:"title" binding:"required"`
	Content  string       `json:"content"`
	Type     DocumentType `json:"type" binding:"required"`
	FolderID *uint        `json:"folder_id"`
	Tags     string       `json:"tags"`
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Title    string       `json:"title"`
	Content  string       `json:"content"`
	Type     DocumentType `json:"type"`
	FolderID *uint        `json:"folder_id"`
	Status   *DocumentStatus `json:"status"`
	Tags     string       `json:"tags"`
}

// DocumentListRequest 文档列表请求
type DocumentListRequest struct {
	ListRequest
	FolderID *uint        `form:"folder_id" json:"folder_id"`
	Type     DocumentType `form:"type" json:"type"`
	Status   DocumentStatus `form:"status" json:"status"`
	Keyword  string       `form:"keyword" json:"keyword"`
}

// ShareDocumentRequest 分享文档请求
type ShareDocumentRequest struct {
	ExpiryDays int  `json:"expiry_days"` // 过期天数，0表示永不过期
	IsPublic   bool `json:"is_public"`   // 是否公开分享
}

// CreateFolderRequest 创建文件夹请求
type CreateFolderRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint  `json:"parent_id"` // 父文件夹ID，null表示根目录
}

// UpdateFolderRequest 更新文件夹请求
type UpdateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

// MoveFolderRequest 移动文件夹请求
type MoveFolderRequest struct {
	ParentID *uint `json:"parent_id"` // 新的父文件夹ID，null表示移至根目录
}

// SearchDocumentsRequest 搜索文档请求
type SearchDocumentsRequest struct {
	ListRequest
	Keyword  string       `form:"keyword" binding:"required" json:"keyword"`
	FolderID *uint        `form:"folder_id" json:"folder_id"`
	Type     DocumentType `form:"type" json:"type"`
}

// SearchAllRequest 全局搜索请求
type SearchAllRequest struct {
	ListRequest
	Keyword string `form:"keyword" binding:"required" json:"keyword"`
}

// CreateActivityRequest 创建活动记录请求
type CreateActivityRequest struct {
	Type        string `json:"type" binding:"required"` // 活动类型：create, update, delete, share, etc.
	ObjectType  string `json:"object_type" binding:"required"` // 对象类型：document, folder
	ObjectID    uint   `json:"object_id" binding:"required"`   // 对象ID
	Description string `json:"description"`                   // 活动描述
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}
