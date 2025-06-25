package model

// DocumentType 文档类型
type DocumentType string

const (
	DocumentTypeWord    DocumentType = "word"    // 文档
	DocumentTypeExcel   DocumentType = "excel"   // 表格
	DocumentTypePPT     DocumentType = "ppt"     // 演示文稿
	DocumentTypeMindMap DocumentType = "mindmap" // 思维导图
	DocumentTypeNote    DocumentType = "note"    // 个人随笔
)

// DocumentStatus 文档状态
type DocumentStatus string

const (
	DocumentStatusDraft     DocumentStatus = "draft"     // 草稿
	DocumentStatusPublished DocumentStatus = "published" // 已发布
	DocumentStatusArchived  DocumentStatus = "archived"  // 已归档
)

// ActivityType 活动类型
type ActivityType string

const (
	ActivityTypeCreate ActivityType = "create" // 创建
	ActivityTypeUpdate ActivityType = "update" // 更新
	ActivityTypeDelete ActivityType = "delete" // 删除
	ActivityTypeShare  ActivityType = "share"  // 分享
	ActivityTypeView   ActivityType = "view"   // 查看
)

// ObjectType 对象类型
type ObjectType string

const (
	ObjectTypeDocument ObjectType = "document" // 文档
	ObjectTypeFolder   ObjectType = "folder"   // 文件夹
)

// SortOrder 排序顺序
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"  // 升序
	SortOrderDesc SortOrder = "desc" // 降序
)
