package service

import (
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/repository"

	"go.uber.org/zap"
)

type ActivityService interface {
	Create(userID uint, req *model.CreateActivityRequest, ipAddress, userAgent string) error
	List(userID uint, req *model.ActivityListRequest) ([]model.ActivityResponse, int64, error)
}

type activityService struct {
	activityRepo repository.ActivityRepository
	logger       *zap.Logger
}

func NewActivityService(activityRepo repository.ActivityRepository, logger *zap.Logger) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
		logger:       logger,
	}
}

func (s *activityService) Create(userID uint, req *model.CreateActivityRequest, ipAddress, userAgent string) error {
	activity := &model.Activity{
		UserID:       userID,
		Type:         req.Type,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
		Description:  req.Description,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	err := s.activityRepo.Create(activity)
	if err != nil {
		return err
	}

	s.logger.Debug("Activity created", 
		zap.Uint("user_id", userID),
		zap.String("type", string(req.Type)),
		zap.String("resource_type", string(req.ResourceType)),
		zap.Uint("resource_id", req.ResourceID))

	return nil
}

func (s *activityService) List(userID uint, req *model.ActivityListRequest) ([]model.ActivityResponse, int64, error) {
	activities, total, err := s.activityRepo.List(userID, req)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.ActivityResponse
	for _, activity := range activities {
		responses = append(responses, model.ActivityResponse{
			ID:           activity.ID,
			Type:         activity.Type,
			ResourceType: activity.ResourceType,
			ResourceID:   activity.ResourceID,
			ResourceName: activity.ResourceName,
			Description:  activity.Description,
			CreatedAt:    activity.CreatedAt,
		})
	}

	return responses, total, nil
}
