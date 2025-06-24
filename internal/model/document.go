package model

import (
	"time"
	"gorm.io/gorm"
)

// DocumentType 文档类型枚举
type DocumentType string

const (
	DocumentTypeWord      DocumentType = "word"      // Word文档
	DocumentTypeExcel     DocumentType = "excel"     // Excel表格
	DocumentTypeMindMap   DocumentType = "mindmap"   // 思维导图
	DocumentTypeNote      DocumentType = "note"      // 个人随笔
	DocumentTypeAIDraft   DocumentType = "ai_draft"  // AI快速起草
	DocumentTypeImported  DocumentType = "imported"  // 导入文件
)

// DocumentStatus 文档状态
type DocumentStatus int

const (
	DocumentStatusDraft     DocumentStatus = 1 // 草稿
	DocumentStatusPublished DocumentStatus = 2 // 已发布
	DocumentStatusArchived  DocumentStatus = 3 // 已归档
)

type Document struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Title       string           `json:"title" gorm:"size:255;not null"`
	Content     string           `json:"content" gorm:"type:longtext"`
	Type        DocumentType     `json:"type" gorm:"size:20;not null"`
	Status      DocumentStatus   `json:"status" gorm:"default:1"`
	UserID      uint             `json:"user_id" gorm:"not null;index"`
	FolderID    *uint            `json:"folder_id" gorm:"index"`
	Tags        string           `json:"tags" gorm:"size:500"`        // 标签，逗号分隔
	Size        int64            `json:"size" gorm:"default:0"`       // 文件大小(字节)
	ViewCount   int              `json:"view_count" gorm:"default:0"` // 查看次数
	IsShared    bool             `json:"is_shared" gorm:"default:false"`
	ShareToken  string           `json:"share_token" gorm:"size:32;index"`
	ShareExpiry *time.Time       `json:"share_expiry"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `json:"-" gorm:"index"`
	
	// 关联
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Folder *Folder `json:"folder,omitempty" gorm:"foreignKey:FolderID"`
}

type Folder struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	ParentID  *uint          `json:"parent_id" gorm:"index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	User     User       `json:"user" gorm:"foreignKey:UserID"`
	Parent   *Folder    `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Folder   `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Documents []Document `json:"documents,omitempty" gorm:"foreignKey:FolderID"`
}

// 请求和响应结构
type CreateDocumentRequest struct {
	Title    string       `json:"title" binding:"required,max=255"`
	Content  string       `json:"content"`
	Type     DocumentType `json:"type" binding:"required"`
	FolderID *uint        `json:"folder_id"`
	Tags     string       `json:"tags"`
}

type UpdateDocumentRequest struct {
	Title    string            `json:"title" binding:"max=255"`
	Content  string            `json:"content"`
	Status   *DocumentStatus   `json:"status"`
	FolderID *uint             `json:"folder_id"`
	Tags     string            `json:"tags"`
}

type DocumentListRequest struct {
	Page     int          `form:"page" binding:"min=1"`
	PageSize int          `form:"page_size" binding:"min=1,max=100"`
	Type     DocumentType `form:"type"`
	Status   DocumentStatus `form:"status"`
	FolderID *uint        `form:"folder_id"`
	Keyword  string       `form:"keyword"`
}

type DocumentResponse struct {
	ID          uint           `json:"id"`
	Title       string         `json:"title"`
	Type        DocumentType   `json:"type"`
	Status      DocumentStatus `json:"status"`
	FolderID    *uint          `json:"folder_id"`
	Tags        string         `json:"tags"`
	Size        int64          `json:"size"`
	ViewCount   int            `json:"view_count"`
	IsShared    bool           `json:"is_shared"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type DocumentDetailResponse struct {
	DocumentResponse
	Content string `json:"content"`
}

type ShareDocumentRequest struct {
	ExpiryHours int `json:"expiry_hours" binding:"min=1,max=8760"` // 最长1年
}

// 文件夹请求和响应结构
type CreateFolderRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateFolderRequest struct {
	Name string `json:"name" binding:"max=100"`
}

type FolderResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FolderTreeResponse struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	ParentID  *uint                `json:"parent_id"`
	Children  []FolderTreeResponse `json:"children"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type MoveFolderRequest struct {
	NewParentID uint `json:"new_parent_id"`
} 