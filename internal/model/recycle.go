package model

import (
	"time"
	"gorm.io/gorm"
)

type RecycleItem struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	ResourceType ResourceType   `json:"resource_type" gorm:"size:20;not null"`
	ResourceID   uint           `json:"resource_id" gorm:"not null"`
	ResourceName string         `json:"resource_name" gorm:"size:255;not null"`
	OriginalPath string         `json:"original_path" gorm:"size:500"` // 原始路径
	DeleteReason string         `json:"delete_reason" gorm:"size:255"`
	AutoDelete   *time.Time     `json:"auto_delete"` // 自动删除时间
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// 请求和响应结构
type RecycleListRequest struct {
	Page         int          `form:"page" binding:"min=1"`
	PageSize     int          `form:"page_size" binding:"min=1,max=100"`
	ResourceType ResourceType `form:"resource_type"`
	Keyword      string       `form:"keyword"`
}

type RecycleResponse struct {
	ID           uint         `json:"id"`
	ResourceType ResourceType `json:"resource_type"`
	ResourceID   uint         `json:"resource_id"`
	ResourceName string       `json:"resource_name"`
	OriginalPath string       `json:"original_path"`
	DeleteReason string       `json:"delete_reason"`
	AutoDelete   *time.Time   `json:"auto_delete"`
	CreatedAt    time.Time    `json:"created_at"`
}

type RestoreRequest struct {
	FolderID *uint `json:"folder_id"` // 可选，指定恢复到的文件夹
}

type DeleteBatchRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
