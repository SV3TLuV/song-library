package service

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"log/slog"
	"song-library-api/src/cmd/api/internal/model"
	"song-library-api/src/cmd/api/internal/repository"
)

var _ GroupService = (*groupService)(nil)

type groupService struct {
	groupRepo repository.GroupRepository
	trManager *manager.Manager
	logger    *slog.Logger
}

func NewGroupService(
	groupRepo repository.GroupRepository,
	trManager *manager.Manager,
	logger *slog.Logger) *groupService {
	return &groupService{
		groupRepo: groupRepo,
		trManager: trManager,
		logger:    logger,
	}
}

func (s *groupService) GetByName(ctx context.Context, group string) (*model.Group, error) {
	groupDB, err := s.groupRepo.GetByName(ctx, group)
	if err != nil {
		return nil, err
	}

	s.logger.Info("get group", "group", group)

	return groupDB, nil
}
