package service

import (
	"context"

	"github.com/Sacs616/streaming-app/services/video/internal/domain"
	"github.com/Sacs616/streaming-app/services/video/internal/repository"
)

type VideoService struct {
	videoRepo repository.VideoRepository
}

func NewVideoService(videoRepo repository.VideoRepository) *VideoService {
	return &VideoService{
		videoRepo: videoRepo,
	}
}

type CreateVideoInput struct {
	Title       string
	Description string
	UploadedBy  uint64
	CategoryID  uint64
}

func (s *VideoService) CreateVideo(ctx context.Context, input *CreateVideoInput) (*domain.Video, error) {
	video := &domain.Video{
		Title:       input.Title,
		Description: input.Description,
		UploadedBy:  input.UploadedBy,
		CategoryID:  input.CategoryID,
		Status:      domain.StatusPending,
		ViewCount:   0,
	}

	if err := s.videoRepo.Create(ctx, video); err != nil {
		return nil, err
	}

	return video, nil
}

func (s *VideoService) GetVideo(ctx context.Context, id uint64) (*domain.Video, error) {
	return s.videoRepo.FindByID(ctx, id)
}

type ListVideosInput struct {
	Page       int
	PageSize   int
	CategoryID *uint64
	Status     *domain.VideoStatus
}

func (s *VideoService) ListVideos(ctx context.Context, input *ListVideosInput) ([]*domain.Video, int64, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 20
	}
	if input.PageSize > 100 {
		input.PageSize = 100
	}

	offset := (input.Page - 1) * input.PageSize
	return s.videoRepo.FindAll(ctx, input.PageSize, offset, input.CategoryID, input.Status)
}

type UpdateVideoInput struct {
	ID           uint64
	Title        string
	Description  string
	Status       domain.VideoStatus
	ThumbnailURL string
	HLSManifest  string
	Duration     int
}

func (s *VideoService) UpdateVideo(ctx context.Context, input *UpdateVideoInput) (*domain.Video, error) {
	video, err := s.videoRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		video.Title = input.Title
	}
	if input.Description != "" {
		video.Description = input.Description
	}
	if input.Status != "" {
		video.Status = input.Status
	}
	if input.ThumbnailURL != "" {
		video.ThumbnailURL = input.ThumbnailURL
	}
	if input.HLSManifest != "" {
		video.HLSManifest = input.HLSManifest
	}
	if input.Duration > 0 {
		video.Duration = input.Duration
	}

	if err := s.videoRepo.Update(ctx, video); err != nil {
		return nil, err
	}

	return video, nil
}

func (s *VideoService) DeleteVideo(ctx context.Context, id uint64) error {
	return s.videoRepo.Delete(ctx, id)
}

func (s *VideoService) IncrementViewCount(ctx context.Context, id uint64) error {
	return s.videoRepo.IncrementViewCount(ctx, id)
}
