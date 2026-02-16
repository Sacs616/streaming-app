package repository

import (
    "context"
    "github.com/Sacs616/streaming-app/services/auth/internal/domain"
    "gorm.io/gorm"
)

type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    FindByID(ctx context.Context, id uint64) (*domain.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}