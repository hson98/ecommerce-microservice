package models

import (
	"github.com/google/uuid"
	userService "github.com/hson98/ecommerce-microservice/src/auth-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" validate:"omitempty"`
	Email     string    `json:"email" validate:"omitempty,lte=60,email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name" validate:"required,lte=30"`
	LastName  string    `json:"last_name" validate:"required,lte=30"`
	Role      string    `json:"role"`
	Base
}

func (User) TableName() string {
	return "users"
}

type UserToken struct {
	User                  *User
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (u *User) userModelToProto() *userService.User {
	return &userService.User{
		Uuid:      u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
