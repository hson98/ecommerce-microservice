package user_usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hson98/ecommerce-microservice/src/auth-service/config"
	"github.com/hson98/ecommerce-microservice/src/auth-service/internal/models"
	user_repository "github.com/hson98/ecommerce-microservice/src/auth-service/internal/user/repository"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/httperrs"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/myjwt"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/utils"
	"time"
)

type UseCase interface {
	Login(ctx context.Context, userLogin *models.User) (*models.UserToken, error)
	Register(ctx context.Context, user *models.User) (*models.UserToken, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type userUC struct {
	userRepo user_repository.Repository
	config   *config.Config
	jwtMaker myjwt.Maker
}

func NewUserUC(userRepo user_repository.Repository, config *config.Config, jwtMaker myjwt.Maker) UseCase {
	return &userUC{
		userRepo: userRepo,
		config:   config,
		jwtMaker: jwtMaker,
	}
}

func (u *userUC) Login(ctx context.Context, userLogin *models.User) (*models.UserToken, error) {
	findUser, err := u.userRepo.FindUser(ctx, map[string]interface{}{"email": userLogin.Email})
	if err != nil {
		return nil, errors.New(httperrs.ErrEmailOrPasswordInvalid)
	}
	errCompare := utils.CheckPassword(findUser.Password, userLogin.Password)
	if errCompare != nil {
		return nil, errors.New(httperrs.ErrEmailOrPasswordInvalid)
	}
	//create access token
	accessToken, accessPayload, err := CreateToken(ctx, u, findUser.ID, u.config.Server.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	//create refresh token
	refreshToken, refreshPayload, err := CreateToken(ctx, u, findUser.ID, u.config.Server.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &models.UserToken{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  findUser,
	}, nil
}
func CreateToken(c context.Context, u *userUC, idUser uuid.UUID, duration time.Duration) (string, *myjwt.Payload, error) {
	//Create token
	token, payload, err := u.jwtMaker.CreateToken(idUser, duration)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}
func (u *userUC) Register(ctx context.Context, dataUser *models.User) (*models.UserToken, error) {
	user, _ := u.userRepo.FindUser(ctx, map[string]interface{}{"email": dataUser.Email})
	if user != nil {
		return nil, errors.New(httperrs.ErrEmailExisted)
	}

	hashedPassword, err := utils.HashPassword(dataUser.Password)
	dataUser.Password = hashedPassword

	if err != nil {
		return nil, err
	}
	userCreated, err := u.userRepo.CreateUser(ctx, dataUser)
	if err != nil {
		return nil, err
	}
	//create access token
	accessToken, accessPayload, err := CreateToken(ctx, u, userCreated.ID, u.config.Server.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	//create refresh token
	refreshToken, refreshPayload, err := CreateToken(ctx, u, userCreated.ID, u.config.Server.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &models.UserToken{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  userCreated,
	}, nil
}

func (u *userUC) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return u.userRepo.FindUser(ctx, map[string]interface{}{"id": userID.String()})
}
