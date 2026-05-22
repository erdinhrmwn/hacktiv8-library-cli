package controller

import (
	"context"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type ActivityController struct {
	activityService *service.ActivityService
}

func NewActivityController(activityService *service.ActivityService) *ActivityController {
	return &ActivityController{activityService: activityService}
}

func (c *ActivityController) GetAll(ctx context.Context) ([]model.ActivityLog, error) {
	return c.activityService.GetAll(ctx)
}

func (c *ActivityController) Log(ctx context.Context, key, description string) error {
	return c.activityService.Log(ctx, key, description)
}
