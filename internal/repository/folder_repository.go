package repository

import (
	"wz-wenzhan-backend/internal/model"
	"gorm.io/gorm"
)

type FolderRepository interface {
	Create(folder *model.Folder) error
	GetByID(id uint) (*model.Folder, error)
	GetByIDAndUserID(id, userID uint) (*model.Folder, error)
	Update(folder *model.Folder) error
	Delete(id, userID uint) error
	GetUserFolders(userID uint) ([]model.Folder, error)
	GetFolderTree(userID uint) ([]model.Folder, error)
	GetSubFolders(parentID uint, userID uint) ([]model.Folder, error)
	CheckFolderExists(name string, parentID *uint, userID uint) (bool, error)
	MoveFolderToParent(folderID, newParentID uint, userID uint) error
}

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{db: db}
}

func (r *folderRepository) Create(folder *model.Folder) error {
	return r.db.Create(folder).Error
}

func (r *folderRepository) GetByID(id uint) (*model.Folder, error) {
	var folder model.Folder
	err := r.db.Preload("User").Preload("Parent").First(&folder, id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *folderRepository) GetByIDAndUserID(id, userID uint) (*model.Folder, error) {
	var folder model.Folder
	err := r.db.Preload("User").Preload("Parent").
		Where("id = ? AND user_id = ?", id, userID).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *folderRepository) Update(folder *model.Folder) error {
	return r.db.Save(folder).Error
}

func (r *folderRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Folder{}).Error
}

func (r *folderRepository) GetUserFolders(userID uint) ([]model.Folder, error) {
	var folders []model.Folder
	err := r.db.Where("user_id = ?", userID).
		Order("parent_id ASC, name ASC").Find(&folders).Error
	return folders, err
}

func (r *folderRepository) GetFolderTree(userID uint) ([]model.Folder, error) {
	var folders []model.Folder
	err := r.db.Preload("Children").
		Where("user_id = ? AND parent_id IS NULL", userID).
		Order("name ASC").Find(&folders).Error
	
	// 递归加载子文件夹
	for i := range folders {
		r.loadChildren(&folders[i])
	}
	
	return folders, err
}

func (r *folderRepository) loadChildren(folder *model.Folder) {
	var children []model.Folder
	r.db.Where("parent_id = ?", folder.ID).
		Order("name ASC").Find(&children)
	
	for i := range children {
		r.loadChildren(&children[i])
	}
	
	folder.Children = children
}

func (r *folderRepository) GetSubFolders(parentID uint, userID uint) ([]model.Folder, error) {
	var folders []model.Folder
	err := r.db.Where("parent_id = ? AND user_id = ?", parentID, userID).
		Order("name ASC").Find(&folders).Error
	return folders, err
}

func (r *folderRepository) CheckFolderExists(name string, parentID *uint, userID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Folder{}).
		Where("name = ? AND user_id = ?", name, userID)
	
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *folderRepository) MoveFolderToParent(folderID, newParentID uint, userID uint) error {
	return r.db.Model(&model.Folder{}).
		Where("id = ? AND user_id = ?", folderID, userID).
		Update("parent_id", newParentID).Error
}
