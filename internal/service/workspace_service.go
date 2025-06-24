package service

import (
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/repository"

	"go.uber.org/zap"
)

type WorkspaceService interface {
	GetDashboard(userID uint) (*model.DashboardData, error)
	GetStats(userID uint) (*model.WorkspaceStats, error)
}

type workspaceService struct {
	workspaceRepo repository.WorkspaceRepository
	logger        *zap.Logger
}

func NewWorkspaceService(workspaceRepo repository.WorkspaceRepository, logger *zap.Logger) WorkspaceService {
	return &workspaceService{
		workspaceRepo: workspaceRepo,
		logger:        logger,
	}
}

func (s *workspaceService) GetDashboard(userID uint) (*model.DashboardData, error) {
	data, err := s.workspaceRepo.GetDashboardData(userID)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Dashboard data retrieved", zap.Uint("user_id", userID))
	return data, nil
}

func (s *workspaceService) GetStats(userID uint) (*model.WorkspaceStats, error) {
	stats, err := s.workspaceRepo.GetStats(userID)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Workspace stats retrieved", zap.Uint("user_id", userID))
	return stats, nil
}
