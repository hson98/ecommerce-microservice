package user_repository

import (
	"context"
	"errors"
	"github.com/hson98/ecommerce-microservice/src/auth-service/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(context context.Context, user *models.User) (*models.User, error)
	FindUser(context context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func (a *userRepo) FindUser(context context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error) {
	var user models.User
	db := a.db.Table(models.User{}.TableName())
	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}
	if err := db.Where(conditions).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (a *userRepo) CreateUser(context context.Context, user *models.User) (*models.User, error) {
	db := a.db.Begin()
	var userInsert models.User
	if err := db.Table(user.TableName()).Create(&user).Scan(&userInsert).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, err
	}

	return &userInsert, nil
}

func NewUserPgRepo(db *gorm.DB) Repository {
	return &userRepo{db: db}
}
