package service

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"song-library-api/src/internal/model"
	"song-library-api/src/internal/repository"
)

var _ GroupService = (*groupService)(nil)

type groupService struct {
	groupRepo repository.GroupRepository
	trManager *manager.Manager
}

func NewGroupService(
	groupRepo repository.GroupRepository,
	trManager *manager.Manager) *groupService {
	return &groupService{
		groupRepo: groupRepo,
		trManager: trManager,
	}
}

func (s *groupService) GetByName(ctx context.Context, group string) (*model.Group, error) {
	return s.groupRepo.GetByName(ctx, group)
}
