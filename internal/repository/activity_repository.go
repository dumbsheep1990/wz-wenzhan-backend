package repository

import (
	"time"
	"wz-wenzhan-backend/internal/model"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(activity *model.Activity) error
	List(userID uint, req *model.ActivityListRequest) ([]model.Activity, int64, error)
	GetRecentActivities(userID uint, limit int) ([]model.Activity, error)
	CountByUserID(userID uint) (int64, error)
	CountByUserIDAndDateRange(userID uint, start, end time.Time) (int64, error)
	GetActivitiesByDay(userID uint, days int) ([]model.ActivityByDay, error)
	DeleteOldActivities(days int) error
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(activity *model.Activity) error {
	return r.db.Create(activity).Error
}

func (r *activityRepository) List(userID uint, req *model.ActivityListRequest) ([]model.Activity, int64, error) {
	var activities []model.Activity
	var total int64
	
	query := r.db.Model(&model.Activity{}).Where("user_id = ?", userID)
	
	// 添加过滤条件
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.ResourceType != "" {
		query = query.Where("resource_type = ?", req.ResourceType)
	}
	if req.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}
	if req.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			endTime = endTime.Add(24 * time.Hour - time.Second) // 包含整天
			query = query.Where("created_at <= ?", endTime)
		}
	}
	
	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Offset(offset).Limit(req.PageSize).
		Order("created_at DESC").Find(&activities).Error
	
	return activities, total, err
}

func (r *activityRepository) GetRecentActivities(userID uint, limit int) ([]model.Activity, error) {
	var activities []model.Activity
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Find(&activities).Error
	return activities, err
}

func (r *activityRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Activity{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *activityRepository) CountByUserIDAndDateRange(userID uint, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&model.Activity{}).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, start, end).
		Count(&count).Error
	return count, err
}

func (r *activityRepository) GetActivitiesByDay(userID uint, days int) ([]model.ActivityByDay, error) {
	var results []model.ActivityByDay
	
	err := r.db.Model(&model.Activity{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, time.Now().AddDate(0, 0, -days)).
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&results).Error
	
	return results, err
}

func (r *activityRepository) DeleteOldActivities(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoffTime).Delete(&model.Activity{}).Error
}
