package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hson98/ecommerce-microservice/src/auth-service/models"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/grpcerrs"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/myjwt"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/utils"
	userService "github.com/hson98/ecommerce-microservice/src/auth-service/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func (u *usersService) Register(ctx context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	user, err := u.registerReqToUserModel(r)
	if err != nil {
		u.logger.Errorf("registerReqToUserModel: %v", err)
		return nil, status.Errorf(grpcerrs.ParseGRPCErrStatusCode(err), "registerReqToUserModel: %v", err)
	}
	if err := utils.ValidateStruct(ctx, user); err != nil {
		u.logger.Errorf("ValidateStruct: %v", err)
		return nil, status.Errorf(grpcerrs.ParseGRPCErrStatusCode(err), "ValidateStruct: %v", err)
	}
	createdUser, err := u.userUC.Register(ctx, user)
	if err != nil {
		u.logger.Errorf("userUC.Register: %v", err)
		return nil, status.Errorf(grpcerrs.ParseGRPCErrStatusCode(err), "Register: %v", err)
	}
	return &userService.RegisterResponse{User: u.userModelToProto(createdUser.User),
		AccessToken:           createdUser.AccessToken,
		AccessTokenExpiresAt:  timestamppb.New(createdUser.AccessTokenExpiresAt),
		RefreshToken:          createdUser.RefreshToken,
		RefreshTokenExpiresAt: timestamppb.New(createdUser.RefreshTokenExpiresAt),
	}, nil
}
func (u *usersService) Login(ctx context.Context, r *userService.LoginRequest) (*userService.LoginResponse, error) {
	userToken, err := u.userUC.Login(ctx, &models.User{Email: r.GetEmail(), Password: r.GetPassword()})
	if err != nil {
		u.logger.Errorf("userUC.Login: %v", err)
		return nil, status.Errorf(grpcerrs.ParseGRPCErrStatusCode(err), "Login: %v", err)
	}
	return &userService.LoginResponse{User: u.userModelToProto(userToken.User),
		AccessToken:           userToken.AccessToken,
		AccessTokenExpiresAt:  timestamppb.New(userToken.AccessTokenExpiresAt),
		RefreshToken:          userToken.RefreshToken,
		RefreshTokenExpiresAt: timestamppb.New(userToken.RefreshTokenExpiresAt),
	}, nil
}

func (u *usersService) GetUserByID(ctx context.Context, r *userService.GetUserByIdRequest) (*userService.GetUserByIdResponse, error) {
	userUUID, err := uuid.Parse(r.GetId())
	if err != nil {
		u.logger.Errorf("uuid.Parse: %v", err)
		return nil, status.Errorf(grpcerrs.ParseGRPCErrStatusCode(err), "uuid.Parse: %v", err)
	}

	user, err := u.userUC.GetUserByID(ctx, userUUID)
	if err != nil {
		u.logger.Errorf("userUC.GetUserByID: %v", err)
		return nil, status.Errorf(grpcerrs.ParseGRPCErrStatusCode(err), "userUC.GetUserByID: %v", err)
	}

	return &userService.GetUserByIdResponse{User: u.userModelToProto(user)}, nil
}

func (u *usersService) GetMe(ctx context.Context, r *userService.GetMeRequest) (*userService.GetMeResponse, error) {
	payload, err := u.getTokenFromCtx(ctx)
	if err != nil {
		u.logger.Errorf("getTokenFromCtx: %v", err)
		return nil, err
	}
	user, err := u.userUC.GetUserByID(ctx, payload.UserID)
	if err != nil {
		u.logger.Errorf("userUC.GetUserByID: %v", err)
		return nil, status.Errorf(grpcerrs.ParseGRPCErrStatusCode(err), "userUC.GetUserByID: %v", err)
	}
	return &userService.GetMeResponse{User: u.userModelToProto(user)}, nil
}

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (u *usersService) getTokenFromCtx(ctx context.Context) (*myjwt.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := u.jwtMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}

func (u *usersService) registerReqToUserModel(r *userService.RegisterRequest) (*models.User, error) {
	candidate := &models.User{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Role:      r.GetRole(),
		Password:  r.GetPassword(),
	}
	//Validate
	return candidate, nil
}
func (u *usersService) userModelToProto(user *models.User) *userService.User {
	userProto := &userService.User{
		Uuid:      user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
	return userProto
}
