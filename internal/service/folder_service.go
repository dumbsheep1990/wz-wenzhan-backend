package service

import (
	"errors"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FolderService interface {
	Create(userID uint, req *model.CreateFolderRequest) (*model.Folder, error)
	GetByID(id, userID uint) (*model.Folder, error)
	Update(id, userID uint, req *model.UpdateFolderRequest) error
	Delete(id, userID uint) error
	GetFolderTree(userID uint) ([]model.FolderTreeResponse, error)
	GetSubFolders(parentID, userID uint) ([]model.FolderResponse, error)
	MoveFolder(id, newParentID, userID uint) error
}

type folderService struct {
	folderRepo   repository.FolderRepository
	documentRepo repository.DocumentRepository
	activityRepo repository.ActivityRepository
	logger       *zap.Logger
}

func NewFolderService(
	folderRepo repository.FolderRepository,
	documentRepo repository.DocumentRepository,
	activityRepo repository.ActivityRepository,
	logger *zap.Logger) FolderService {
	return &folderService{
		folderRepo:   folderRepo,
		documentRepo: documentRepo,
		activityRepo: activityRepo,
		logger:       logger,
	}
}

func (s *folderService) Create(userID uint, req *model.CreateFolderRequest) (*model.Folder, error) {
	// 检查文件夹名称是否已存在
	exists, err := s.folderRepo.CheckFolderExists(req.Name, req.ParentID, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("文件夹名称已存在")
	}

	// 如果有父文件夹，检查父文件夹是否存在且属于当前用户
	if req.ParentID != nil {
		_, err := s.folderRepo.GetByIDAndUserID(*req.ParentID, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父文件夹不存在")
			}
			return nil, err
		}
	}

	folder := &model.Folder{
		Name:     req.Name,
		UserID:   userID,
		ParentID: req.ParentID,
	}

	err = s.folderRepo.Create(folder)
	if err != nil {
		return nil, err
	}

	// 记录活动
	go func() {
		activity := &model.Activity{
			UserID:       userID,
			Type:         model.ActivityTypeCreate,
			ResourceType: model.ResourceTypeFolder,
			ResourceID:   folder.ID,
			ResourceName: folder.Name,
			Description:  "创建文件夹",
		}
		s.activityRepo.Create(activity)
	}()

	s.logger.Info("Folder created",
		zap.Uint("user_id", userID),
		zap.Uint("folder_id", folder.ID),
		zap.String("name", folder.Name))

	return folder, nil
}

func (s *folderService) GetByID(id, userID uint) (*model.Folder, error) {
	return s.folderRepo.GetByIDAndUserID(id, userID)
}

func (s *folderService) Update(id, userID uint, req *model.UpdateFolderRequest) error {
	folder, err := s.folderRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	oldName := folder.Name

	// 如果要修改名称，检查新名称是否已存在
	if req.Name != "" && req.Name != folder.Name {
		exists, err := s.folderRepo.CheckFolderExists(req.Name, folder.ParentID, userID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("文件夹名称已存在")
		}
		folder.Name = req.Name
	}

	err = s.folderRepo.Update(folder)
	if err != nil {
		return err
	}

	// 记录活动
	go func() {
		activity := &model.Activity{
			UserID:       userID,
			Type:         model.ActivityTypeUpdate,
			ResourceType: model.ResourceTypeFolder,
			ResourceID:   folder.ID,
			ResourceName: folder.Name,
			Description:  "更新文件夹：" + oldName + " -> " + folder.Name,
		}
		s.activityRepo.Create(activity)
	}()

	s.logger.Info("Folder updated",
		zap.Uint("user_id", userID),
		zap.Uint("folder_id", id))

	return nil
}

func (s *folderService) Delete(id, userID uint) error {
	folder, err := s.folderRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	// 检查是否有子文件夹
	subFolders, err := s.folderRepo.GetSubFolders(id, userID)
	if err != nil {
		return err
	}
	if len(subFolders) > 0 {
		return errors.New("无法删除包含子文件夹的文件夹")
	}

	err = s.folderRepo.Delete(id, userID)
	if err != nil {
		return err
	}

	// 记录活动
	go func() {
		activity := &model.Activity{
			UserID:       userID,
			Type:         model.ActivityTypeDelete,
			ResourceType: model.ResourceTypeFolder,
			ResourceID:   folder.ID,
			ResourceName: folder.Name,
			Description:  "删除文件夹",
		}
		s.activityRepo.Create(activity)
	}()

	s.logger.Info("Folder deleted",
		zap.Uint("user_id", userID),
		zap.Uint("folder_id", id))

	return nil
}

func (s *folderService) GetFolderTree(userID uint) ([]model.FolderTreeResponse, error) {
	folders, err := s.folderRepo.GetFolderTree(userID)
	if err != nil {
		return nil, err
	}

	var responses []model.FolderTreeResponse
	for _, folder := range folders {
		responses = append(responses, s.buildFolderTreeResponse(folder))
	}

	return responses, nil
}

func (s *folderService) buildFolderTreeResponse(folder model.Folder) model.FolderTreeResponse {
	response := model.FolderTreeResponse{
		ID:        folder.ID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		CreatedAt: folder.CreatedAt,
		UpdatedAt: folder.UpdatedAt,
	}

	for _, child := range folder.Children {
		response.Children = append(response.Children, s.buildFolderTreeResponse(child))
	}

	return response
}

func (s *folderService) GetSubFolders(parentID, userID uint) ([]model.FolderResponse, error) {
	folders, err := s.folderRepo.GetSubFolders(parentID, userID)
	if err != nil {
		return nil, err
	}

	var responses []model.FolderResponse
	for _, folder := range folders {
		responses = append(responses, model.FolderResponse{
			ID:        folder.ID,
			Name:      folder.Name,
			ParentID:  folder.ParentID,
			CreatedAt: folder.CreatedAt,
			UpdatedAt: folder.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *folderService) MoveFolder(id, newParentID, userID uint) error {
	// 检查文件夹是否存在且属于当前用户
	_, err := s.folderRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	// 检查新父文件夹是否存在且属于当前用户
	if newParentID != 0 {
		_, err := s.folderRepo.GetByIDAndUserID(newParentID, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("目标文件夹不存在")
			}
			return err
		}
	}

	// 移动文件夹
	err = s.folderRepo.MoveFolderToParent(id, newParentID, userID)
	if err != nil {
		return err
	}

	s.logger.Info("Folder moved",
		zap.Uint("user_id", userID),
		zap.Uint("folder_id", id),
		zap.Uint("new_parent_id", newParentID))

	return nil
}
