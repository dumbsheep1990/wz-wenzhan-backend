package repository

import (
	"wz-wenzhan-backend/internal/model"
	"gorm.io/gorm"
)

type RecycleRepository interface {
	Create(item *model.RecycleItem) error
	List(userID uint, req *model.RecycleListRequest) ([]model.RecycleItem, int64, error)
	GetByID(id uint) (*model.RecycleItem, error)
	GetByIDAndUserID(id, userID uint) (*model.RecycleItem, error)
	Delete(id, userID uint) error
	DeleteBatch(ids []uint, userID uint) error
	CountByUserID(userID uint) (int64, error)
	DeleteExpired() error
}

type recycleRepository struct {
	db *gorm.DB
}

func NewRecycleRepository(db *gorm.DB) RecycleRepository {
	return &recycleRepository{db: db}
}

func (r *recycleRepository) Create(item *model.RecycleItem) error {
	return r.db.Create(item).Error
}

func (r *recycleRepository) List(userID uint, req *model.RecycleListRequest) ([]model.RecycleItem, int64, error) {
	var items []model.RecycleItem
	var total int64
	
	query := r.db.Model(&model.RecycleItem{}).Where("user_id = ?", userID)
	
	// 添加过滤条件
	if req.ResourceType != "" {
		query = query.Where("resource_type = ?", req.ResourceType)
	}
	if req.Keyword != "" {
		query = query.Where("resource_name LIKE ?", "%"+req.Keyword+"%")
	}
	
	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Offset(offset).Limit(req.PageSize).
		Order("created_at DESC").Find(&items).Error
	
	return items, total, err
}

func (r *recycleRepository) GetByID(id uint) (*model.RecycleItem, error) {
	var item model.RecycleItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *recycleRepository) GetByIDAndUserID(id, userID uint) (*model.RecycleItem, error) {
	var item model.RecycleItem
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *recycleRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RecycleItem{}).Error
}

func (r *recycleRepository) DeleteBatch(ids []uint, userID uint) error {
	return r.db.Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.RecycleItem{}).Error
}

func (r *recycleRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.RecycleItem{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *recycleRepository) DeleteExpired() error {
	return r.db.Where("auto_delete IS NOT NULL AND auto_delete <= NOW()").Delete(&model.RecycleItem{}).Error
}
