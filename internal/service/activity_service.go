package service

import (
	"context"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type ActivityService struct {
	activityRepository *repository.ActivityRepository
}

func NewActivityService(activityRepository *repository.ActivityRepository) *ActivityService {
	return &ActivityService{activityRepository: activityRepository}
}

func (s *ActivityService) GetAll(ctx context.Context) ([]model.ActivityLog, error) {
	return s.activityRepository.GetAll(ctx)
}

func (s *ActivityService) Log(ctx context.Context, key, description string) error {
	return s.activityRepository.Log(ctx, key, description)
}
