package repository

import (
	"time"
	"wz-wenzhan-backend/internal/model"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(document *model.Document) error
	GetByID(id uint) (*model.Document, error)
	GetByIDAndUserID(id, userID uint) (*model.Document, error)
	Update(document *model.Document) error
	Delete(id, userID uint) error
	List(userID uint, req *model.DocumentListRequest) ([]model.Document, int64, error)
	GetByShareToken(token string) (*model.Document, error)
	UpdateViewCount(id uint) error
	GetRecentDocuments(userID uint, limit int) ([]model.Document, error)
	CountByUserID(userID uint) (int64, error)
	CountByUserIDAndStatus(userID uint, status model.DocumentStatus) (int64, error)
	CountByType(userID uint) (map[string]int64, error)
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(document *model.Document) error {
	return r.db.Create(document).Error
}

func (r *documentRepository) GetByID(id uint) (*model.Document, error) {
	var document model.Document
	err := r.db.Preload("User").Preload("Folder").First(&document, id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) GetByIDAndUserID(id, userID uint) (*model.Document, error) {
	var document model.Document
	err := r.db.Preload("User").Preload("Folder").
		Where("id = ? AND user_id = ?", id, userID).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) Update(document *model.Document) error {
	return r.db.Save(document).Error
}

func (r *documentRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Document{}).Error
}

func (r *documentRepository) List(userID uint, req *model.DocumentListRequest) ([]model.Document, int64, error) {
	var documents []model.Document
	var total int64
	
	query := r.db.Model(&model.Document{}).Where("user_id = ?", userID)
	
	// 添加过滤条件
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.FolderID != nil {
		query = query.Where("folder_id = ?", *req.FolderID)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}
	
	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Preload("Folder").
		Offset(offset).Limit(req.PageSize).
		Order("updated_at DESC").Find(&documents).Error
	
	return documents, total, err
}

func (r *documentRepository) GetByShareToken(token string) (*model.Document, error) {
	var document model.Document
	err := r.db.Preload("User").
		Where("share_token = ? AND is_shared = ? AND (share_expiry IS NULL OR share_expiry > ?)", 
			token, true, time.Now()).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) UpdateViewCount(id uint) error {
	return r.db.Model(&model.Document{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *documentRepository) GetRecentDocuments(userID uint, limit int) ([]model.Document, error) {
	var documents []model.Document
	err := r.db.Where("user_id = ?", userID).
		Order("updated_at DESC").Limit(limit).Find(&documents).Error
	return documents, err
}

func (r *documentRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Document{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *documentRepository) CountByUserIDAndStatus(userID uint, status model.DocumentStatus) (int64, error) {
	var count int64
	err := r.db.Model(&model.Document{}).
		Where("user_id = ? AND status = ?", userID, status).Count(&count).Error
	return count, err
}

func (r *documentRepository) CountByType(userID uint) (map[string]int64, error) {
	var results []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	
	err := r.db.Model(&model.Document{}).
		Select("type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("type").Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	countMap := make(map[string]int64)
	for _, result := range results {
		countMap[result.Type] = result.Count
	}
	
	return countMap, nil
} 