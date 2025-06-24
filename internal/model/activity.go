package model

import (
	"time"
	"gorm.io/gorm"
)

// ActivityType 活动类型
type ActivityType string

const (
	ActivityTypeCreate ActivityType = "create" // 创建
	ActivityTypeUpdate ActivityType = "update" // 更新
	ActivityTypeDelete ActivityType = "delete" // 删除
	ActivityTypeView   ActivityType = "view"   // 查看
	ActivityTypeShare  ActivityType = "share"  // 分享
	ActivityTypeCopy   ActivityType = "copy"   // 复制
	ActivityTypeMove   ActivityType = "move"   // 移动
)

// ResourceType 资源类型
type ResourceType string

const (
	ResourceTypeDocument ResourceType = "document" // 文档
	ResourceTypeFolder   ResourceType = "folder"   // 文件夹
)

type Activity struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	Type         ActivityType   `json:"type" gorm:"size:20;not null"`
	ResourceType ResourceType   `json:"resource_type" gorm:"size:20;not null"`
	ResourceID   uint           `json:"resource_id" gorm:"not null"`
	ResourceName string         `json:"resource_name" gorm:"size:255;not null"`
	Description  string         `json:"description" gorm:"size:500"`
	IPAddress    string         `json:"ip_address" gorm:"size:45"`
	UserAgent    string         `json:"user_agent" gorm:"size:500"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// 请求和响应结构
type CreateActivityRequest struct {
	Type         ActivityType `json:"type" binding:"required"`
	ResourceType ResourceType `json:"resource_type" binding:"required"`
	ResourceID   uint         `json:"resource_id" binding:"required"`
	ResourceName string       `json:"resource_name" binding:"required"`
	Description  string       `json:"description"`
}

type ActivityListRequest struct {
	Page         int          `form:"page" binding:"min=1"`
	PageSize     int          `form:"page_size" binding:"min=1,max=100"`
	Type         ActivityType `form:"type"`
	ResourceType ResourceType `form:"resource_type"`
	StartDate    string       `form:"start_date"` // YYYY-MM-DD
	EndDate      string       `form:"end_date"`   // YYYY-MM-DD
}

type ActivityResponse struct {
	ID           uint         `json:"id"`
	Type         ActivityType `json:"type"`
	ResourceType ResourceType `json:"resource_type"`
	ResourceID   uint         `json:"resource_id"`
	ResourceName string       `json:"resource_name"`
	Description  string       `json:"description"`
	CreatedAt    time.Time    `json:"created_at"`
} 