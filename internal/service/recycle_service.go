package service

import (
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/repository"

	"go.uber.org/zap"
)

type RecycleService interface {
	List(userID uint, req *model.RecycleListRequest) ([]model.RecycleResponse, int64, error)
	Restore(id, userID uint, req *model.RestoreRequest) error
	DeletePermanently(id, userID uint) error
	DeleteBatch(userID uint, req *model.DeleteBatchRequest) error
}

type recycleService struct {
	recycleRepo  repository.RecycleRepository
	documentRepo repository.DocumentRepository
	logger       *zap.Logger
}

func NewRecycleService(recycleRepo repository.RecycleRepository, logger *zap.Logger) RecycleService {
	return &recycleService{
		recycleRepo: recycleRepo,
		logger:      logger,
	}
}

func (s *recycleService) List(userID uint, req *model.RecycleListRequest) ([]model.RecycleResponse, int64, error) {
	items, total, err := s.recycleRepo.List(userID, req)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.RecycleResponse
	for _, item := range items {
		responses = append(responses, model.RecycleResponse{
			ID:           item.ID,
			ResourceType: item.ResourceType,
			ResourceID:   item.ResourceID,
			ResourceName: item.ResourceName,
			OriginalPath: item.OriginalPath,
			DeleteReason: item.DeleteReason,
			AutoDelete:   item.AutoDelete,
			CreatedAt:    item.CreatedAt,
		})
	}

	return responses, total, nil
}

func (s *recycleService) Restore(id, userID uint, req *model.RestoreRequest) error {
	// TODO: 实现恢复逻辑
	// 1. 获取回收站项目
	// 2. 根据资源类型恢复相应的资源
	// 3. 删除回收站记录
	
	s.logger.Info("Item restored from recycle", 
		zap.Uint("user_id", userID),
		zap.Uint("recycle_id", id))

	return nil
}

func (s *recycleService) DeletePermanently(id, userID uint) error {
	err := s.recycleRepo.Delete(id, userID)
	if err != nil {
		return err
	}

	s.logger.Info("Item permanently deleted", 
		zap.Uint("user_id", userID),
		zap.Uint("recycle_id", id))

	return nil
}

func (s *recycleService) DeleteBatch(userID uint, req *model.DeleteBatchRequest) error {
	err := s.recycleRepo.DeleteBatch(req.IDs, userID)
	if err != nil {
		return err
	}

	s.logger.Info("Items batch deleted", 
		zap.Uint("user_id", userID),
		zap.Int("count", len(req.IDs)))

	return nil
}
