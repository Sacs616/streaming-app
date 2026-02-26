package repository

import (
	"context"

	"github.com/Sacs616/streaming-app/services/video/internal/domain"
	"gorm.io/gorm"
)

type VideoRepository interface {
	Create(ctx context.Context, video *domain.Video) error
	FindByID(ctx context.Context, id uint64) (*domain.Video, error)
	FindAll(ctx context.Context, limit, offset int, categoryID *uint64, status *domain.VideoStatus) ([]*domain.Video, int64, error)
	Update(ctx context.Context, video *domain.Video) error
	Delete(ctx context.Context, id uint64) error
	IncrementViewCount(ctx context.Context, id uint64) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(ctx context.Context, video *domain.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) FindByID(ctx context.Context, id uint64) (*domain.Video, error) {
	var video domain.Video
	err := r.db.WithContext(ctx).First(&video, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrVideoNotFound
		}
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindAll(
	ctx context.Context,
	limit, offset int,
	categoryID *uint64,
	status *domain.VideoStatus,
) ([]*domain.Video, int64, error) {
	var videos []*domain.Video
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Video{})

	// Apply filters
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error

	return videos, total, err
}

func (r *videoRepository) Update(ctx context.Context, video *domain.Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

func (r *videoRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Video{}, id).Error
}

func (r *videoRepository) IncrementViewCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Video{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
		Error
}
