package user_usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hson98/ecommerce-microservice/src/auth-service/config"
	user_repository "github.com/hson98/ecommerce-microservice/src/auth-service/internal/user/repository"
	"github.com/hson98/ecommerce-microservice/src/auth-service/models"
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

type authUC struct {
	authRepo user_repository.Repository
	config   *config.Config
	jwtMaker myjwt.Maker
}

func NewAuthUC(authRepo user_repository.Repository, config *config.Config, jwtMaker myjwt.Maker) UseCase {
	return &authUC{
		authRepo: authRepo,
		config:   config,
		jwtMaker: jwtMaker,
	}
}

func (a *authUC) Login(ctx context.Context, userLogin *models.User) (*models.UserToken, error) {
	findUser, err := a.authRepo.FindUser(ctx, map[string]interface{}{"email": userLogin.Email})
	if err != nil {
		return nil, errors.New(httperrs.ErrUsernameOrPasswordInvalid)
	}
	errCompare := utils.CheckPassword(findUser.Password, userLogin.Password)
	if errCompare != nil {
		return nil, errors.New(httperrs.ErrUsernameOrPasswordInvalid)
	}
	//create access token
	accessToken, accessPayload, err := CreateToken(ctx, a, findUser.ID, a.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	//create refresh token
	refreshToken, refreshPayload, err := CreateToken(ctx, a, findUser.ID, a.config.RefreshTokenDuration)
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
func CreateToken(c context.Context, a *authUC, idUser uuid.UUID, duration time.Duration) (string, *myjwt.Payload, error) {
	//Create token
	token, payload, err := a.jwtMaker.CreateToken(idUser, duration)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}
func (a *authUC) Register(ctx context.Context, dataUser *models.User) (*models.UserToken, error) {
	user, _ := a.authRepo.FindUser(ctx, map[string]interface{}{"email": dataUser.Email})
	if user != nil {
		return nil, errors.New(httperrs.ErrEmailExisted)
	}

	hashedPassword, err := utils.HashPassword(dataUser.Password)
	dataUser.Password = hashedPassword

	if err != nil {
		return nil, err
	}
	userCreated, err := a.authRepo.CreateUser(ctx, dataUser)
	if err != nil {
		return nil, err
	}
	//create access token
	accessToken, accessPayload, err := CreateToken(ctx, a, userCreated.ID, a.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	//create refresh token
	refreshToken, refreshPayload, err := CreateToken(ctx, a, userCreated.ID, a.config.RefreshTokenDuration)
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

func (a *authUC) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return a.authRepo.FindUser(ctx, map[string]interface{}{"id": userID.String()})
}
