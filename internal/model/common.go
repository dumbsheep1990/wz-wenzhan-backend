package model

import "time"

// BaseModel 基础模型，包含公共字段
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int `form:"page" binding:"min=1" json:"page"`
	PageSize int `form:"page_size" binding:"min=1,max=100" json:"page_size"`
}

// SetDefaults 设置分页默认值
func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(items interface{}, total int64, page, pageSize int) *PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &PaginationResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Code:      200,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *APIResponse {
	return &APIResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// SearchRequest 搜索请求
type SearchRequest struct {
	PaginationRequest
	Keyword    string `form:"keyword" json:"keyword"`
	StartDate  string `form:"start_date" json:"start_date"` // YYYY-MM-DD
	EndDate    string `form:"end_date" json:"end_date"`     // YYYY-MM-DD
	SortBy     string `form:"sort_by" json:"sort_by"`
	SortOrder  string `form:"sort_order" json:"sort_order"` // asc, desc
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}
