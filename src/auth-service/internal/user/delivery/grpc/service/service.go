package service

import (
	"github.com/hson98/ecommerce-microservice/src/auth-service/config"
	user_usecase "github.com/hson98/ecommerce-microservice/src/auth-service/internal/user/usecase"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/logger"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/myjwt"
)

type usersService struct {
	logger   logger.Logger
	cfg      *config.Config
	userUC   user_usecase.UseCase
	jwtMaker myjwt.Maker
}

func NewAuthServerGRPC(logger logger.Logger, cfg *config.Config, userUC user_usecase.UseCase, jwtMaker myjwt.Maker) *usersService {
	return &usersService{
		logger:   logger,
		cfg:      cfg,
		userUC:   userUC,
		jwtMaker: jwtMaker,
	}
}
