package repository

import (
	"time"
	"wz-wenzhan-backend/internal/model"
	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	GetStats(userID uint) (*model.WorkspaceStats, error)
	GetDashboardData(userID uint) (*model.DashboardData, error)
}

type workspaceRepository struct {
	db           *gorm.DB
	documentRepo DocumentRepository
	activityRepo ActivityRepository
	recycleRepo  RecycleRepository
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{
		db:           db,
		documentRepo: NewDocumentRepository(db),
		activityRepo: NewActivityRepository(db),
		recycleRepo:  NewRecycleRepository(db),
	}
}

func (r *workspaceRepository) GetStats(userID uint) (*model.WorkspaceStats, error) {
	stats := &model.WorkspaceStats{}
	
	// 获取文档统计
	totalDocs, err := r.documentRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	stats.TotalDocuments = totalDocs
	
	draftDocs, err := r.documentRepo.CountByUserIDAndStatus(userID, model.DocumentStatusDraft)
	if err != nil {
		return nil, err
	}
	stats.DraftDocuments = draftDocs
	
	publishedDocs, err := r.documentRepo.CountByUserIDAndStatus(userID, model.DocumentStatusPublished)
	if err != nil {
		return nil, err
	}
	stats.PublishedDocuments = publishedDocs
	
	// 获取活动统计
	totalActivities, err := r.activityRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	stats.TotalActivities = totalActivities
	
	// 今日活动
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour - time.Second)
	
	todayActivities, err := r.activityRepo.CountByUserIDAndDateRange(userID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	stats.TodayActivities = todayActivities
	
	// 本周活动
	weekAgo := today.AddDate(0, 0, -7)
	weekActivities, err := r.activityRepo.CountByUserIDAndDateRange(userID, weekAgo, today)
	if err != nil {
		return nil, err
	}
	stats.WeekActivities = weekActivities
	
	// 回收站数量
	recycleCount, err := r.recycleRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	stats.RecycleItems = recycleCount
	
	return stats, nil
}

func (r *workspaceRepository) GetDashboardData(userID uint) (*model.DashboardData, error) {
	// 获取统计数据
	stats, err := r.GetStats(userID)
	if err != nil {
		return nil, err
	}
	
	// 获取最近文档
	recentDocs, err := r.documentRepo.GetRecentDocuments(userID, 5)
	if err != nil {
		return nil, err
	}
	
	var recentDocuments []model.DocumentResponse
	for _, doc := range recentDocs {
		recentDocuments = append(recentDocuments, model.DocumentResponse{
			ID:        doc.ID,
			Title:     doc.Title,
			Type:      doc.Type,
			Status:    doc.Status,
			FolderID:  doc.FolderID,
			Tags:      doc.Tags,
			Size:      doc.Size,
			ViewCount: doc.ViewCount,
			IsShared:  doc.IsShared,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		})
	}
	
	// 获取最近活动
	recentActs, err := r.activityRepo.GetRecentActivities(userID, 10)
	if err != nil {
		return nil, err
	}
	
	var recentActivities []model.ActivityResponse
	for _, act := range recentActs {
		recentActivities = append(recentActivities, model.ActivityResponse{
			ID:           act.ID,
			Type:         act.Type,
			ResourceType: act.ResourceType,
			ResourceID:   act.ResourceID,
			ResourceName: act.ResourceName,
			Description:  act.Description,
			CreatedAt:    act.CreatedAt,
		})
	}
	
	// 获取文档类型统计
	docsByType, err := r.documentRepo.CountByType(userID)
	if err != nil {
		return nil, err
	}
	
	// 获取最近7天的活动统计
	activitiesByDay, err := r.activityRepo.GetActivitiesByDay(userID, 7)
	if err != nil {
		return nil, err
	}
	
	dashboardData := &model.DashboardData{
		Stats:             *stats,
		RecentDocuments:   recentDocuments,
		RecentActivities:  recentActivities,
		DocumentsByType:   docsByType,
		ActivitiesByDay:   activitiesByDay,
	}
	
	return dashboardData, nil
}
