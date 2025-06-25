package model

import (
	"time"
)

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"`
	User  *UserResponse `json:"user"`
}

// DocumentResponse 文档响应
type DocumentResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Content     string        `json:"content,omitempty"` // 列表时可能不需要返回内容
	Type        DocumentType  `json:"type"`
	Status      DocumentStatus `json:"status"`
	UserID      uint          `json:"user_id"`
	FolderID    *uint         `json:"folder_id"`
	Tags        string        `json:"tags"`
	Size        int64         `json:"size"`
	ViewCount   int           `json:"view_count"`
	IsShared    bool          `json:"is_shared"`
	ShareToken  string        `json:"share_token,omitempty"`
	ShareExpiry *time.Time    `json:"share_expiry,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// DocumentListResponse 文档列表响应
type DocumentListResponse struct {
	Items []*DocumentResponse `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// FolderResponse 文件夹响应
type FolderResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	UserID    uint      `json:"user_id"`
	ItemCount int       `json:"item_count"` // 包含的文档和子文件夹数量
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FolderTreeResponse 文件夹树响应
type FolderTreeResponse struct {
	FolderResponse
	Children []*FolderTreeResponse `json:"children,omitempty"`
}

// FileResponse 文件响应
type FileResponse struct {
	ID        uint      `json:"id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	UserID    uint      `json:"user_id"`
	FolderID  *uint     `json:"folder_id"`
	CreatedAt time.Time `json:"created_at"`
}

// FileURLResponse 文件URL响应
type FileURLResponse struct {
	URL string `json:"url"`
}

// ShareResponse 分享响应
type ShareResponse struct {
	Token    string `json:"token"`
	ShareURL string `json:"share_url"`
}

// ActivityResponse 活动记录响应
type ActivityResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Type        string    `json:"type"`
	ObjectType  string    `json:"object_type"`
	ObjectID    uint      `json:"object_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ActivityListResponse 活动记录列表响应
type ActivityListResponse struct {
	Items []*ActivityResponse `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// RecycleItemResponse 回收站项目响应
type RecycleItemResponse struct {
	ID         uint      `json:"id"`
	ObjectType string    `json:"object_type"` // document, folder
	ObjectID   uint      `json:"object_id"`
	Title      string    `json:"title"`
	UserID     uint      `json:"user_id"`
	DeletedAt  time.Time `json:"deleted_at"`
}

// RecycleListResponse 回收站列表响应
type RecycleListResponse struct {
	Items []*RecycleItemResponse `json:"items"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// SearchResultItem 搜索结果项
type SearchResultItem struct {
	ID         uint      `json:"id"`
	Type       string    `json:"type"` // document, folder
	Title      string    `json:"title"`
	Content    string    `json:"content,omitempty"` // 命中的内容片段
	FolderID   *uint     `json:"folder_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	MatchScore float64   `json:"match_score"`
}

// SearchDocumentsResponse 搜索文档响应
type SearchDocumentsResponse struct {
	Items []*DocumentResponse `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// SearchAllResponse 全局搜索响应
type SearchAllResponse struct {
	Items []*SearchResultItem `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// DashboardResponse 工作台首页响应
type DashboardResponse struct {
	RecentDocuments    []*DocumentResponse `json:"recent_documents"`
	ImportantDocuments []*DocumentResponse `json:"important_documents"`
	Activities         []*ActivityResponse `json:"activities"`
}

// StatsResponse 统计信息响应
type StatsResponse struct {
	DocumentCount    int64            `json:"document_count"`
	FolderCount      int64            `json:"folder_count"`
	StorageUsed      int64            `json:"storage_used"`       // 字节
	StorageLimit     int64            `json:"storage_limit"`      // 字节
	TypeDistribution map[string]int64 `json:"type_distribution"` // 各类型文档数量
}
