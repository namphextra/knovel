package repository

import (
    "context"
    "knovel/internal/models"
    "gorm.io/gorm"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type GormUserRepository struct {
    db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
    return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    return &user, err
} 