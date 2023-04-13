package services

import (
	"context"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/pkg/logging"
)

type GameService struct {
	gameRepo repository.GameRepositoryInterface
	roleRepo repository.RoleRepositoryInterface
	userRepo repository.UserRepositoryInterface
	logger   *logging.Logger
	ctx      context.Context
}

func NewGameService(
	gameRepo repository.GameRepositoryInterface,
	roleRepo repository.RoleRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	logger *logging.Logger,
	ctx context.Context) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		roleRepo: roleRepo,
		userRepo: userRepo,
		logger:   logger,
		ctx:      ctx,
	}
}

func (s *GameService) GetOneById(id int) (*models.Game, error) {
	return s.gameRepo.GetOneById(s.ctx, id)
}

func (s *GameService) AddUserToGame(id int, user *models.User) (*models.Game, error) {
	game, err := s.gameRepo.GetOneById(s.ctx, id)
	if err != nil {
		return nil, err
	}

	defaultRole, err := s.roleRepo.GetOne(s.ctx, "user")
	if err != nil {
		return nil, err
	}

	err1 := s.userRepo.AddUserToGameWithRole(s.ctx, user, game, defaultRole)
	if err1 != nil {
		return nil, err1
	}

	return game, nil
}

func (s *GameService) GetUserGames(user *models.User) ([]repository.UserGameInfo, error) {
	games, err := s.userRepo.GetUserGames(s.ctx, user)
	if err != nil {
		return nil, err
	}
	return games, err
}
