package model

import "time"

// WorkspaceStats 工作台统计数据
type WorkspaceStats struct {
	TotalDocuments   int64 `json:"total_documents"`
	DraftDocuments   int64 `json:"draft_documents"`
	PublishedDocuments int64 `json:"published_documents"`
	TotalActivities  int64 `json:"total_activities"`
	TodayActivities  int64 `json:"today_activities"`
	WeekActivities   int64 `json:"week_activities"`
	RecycleItems     int64 `json:"recycle_items"`
}

// DashboardData 仪表板数据
type DashboardData struct {
	Stats             WorkspaceStats      `json:"stats"`
	RecentDocuments   []DocumentResponse  `json:"recent_documents"`
	RecentActivities  []ActivityResponse  `json:"recent_activities"`
	DocumentsByType   map[string]int64    `json:"documents_by_type"`
	ActivitiesByDay   []ActivityByDay     `json:"activities_by_day"`
}

// ActivityByDay 按天统计的活动数据
type ActivityByDay struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// QuickAction 快速操作
type QuickAction struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	URL   string `json:"url"`
}

// GetQuickActions 获取快速操作列表
func GetQuickActions() []QuickAction {
	return []QuickAction{
		{Type: "word", Title: "Word文档", Icon: "file-word", URL: "/documents/create?type=word"},
		{Type: "excel", Title: "Excel表格", Icon: "file-excel", URL: "/documents/create?type=excel"},
		{Type: "mindmap", Title: "思维导图", Icon: "share-alt", URL: "/documents/create?type=mindmap"},
		{Type: "note", Title: "个人随笔", Icon: "edit", URL: "/documents/create?type=note"},
		{Type: "ai_draft", Title: "AI快速起草", Icon: "robot", URL: "/documents/create?type=ai_draft"},
		{Type: "import", Title: "文件导入", Icon: "upload", URL: "/documents/import"},
	}
} 