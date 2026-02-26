package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Sacs616/streaming-app/services/video/internal/domain"
	"github.com/Sacs616/streaming-app/services/video/internal/service"
	videopb "github.com/Sacs616/streaming-app/services/video/proto"
	commonpb "github.com/Sacs616/streaming-app/shared/proto/common"
)

type VideoServer struct {
	videopb.UnimplementedVideoServiceServer
	videoService *service.VideoService
}

func NewVideoServer(videoService *service.VideoService) *VideoServer {
	return &VideoServer{
		videoService: videoService,
	}
}

func (s *VideoServer) CreateVideo(ctx context.Context, req *videopb.CreateVideoRequest) (*videopb.CreateVideoResponse, error) {
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "category_id is required")
	}

	video, err := s.videoService.CreateVideo(ctx, &service.CreateVideoInput{
		Title:       req.Title,
		Description: req.Description,
		UploadedBy:  req.UploadedBy,
		CategoryID:  req.CategoryId,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create video")
	}

	return &videopb.CreateVideoResponse{
		Video: s.mapToProto(video),
	}, nil
}

func (s *VideoServer) GetVideo(ctx context.Context, req *videopb.GetVideoRequest) (*videopb.GetVideoResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	video, err := s.videoService.GetVideo(ctx, req.Id)
	if err != nil {
		if err == domain.ErrVideoNotFound {
			return nil, status.Error(codes.NotFound, "video not found")
		}
		return nil, status.Error(codes.Internal, "failed to get video")
	}

	// Increment view count asynchronously
	go s.videoService.IncrementViewCount(context.Background(), req.Id)

	return &videopb.GetVideoResponse{
		Video: s.mapToProto(video),
	}, nil
}

func (s *VideoServer) ListVideos(ctx context.Context, req *videopb.ListVideosRequest) (*videopb.ListVideosResponse, error) {
	page := int(req.Pagination.Page)
	pageSize := int(req.Pagination.PageSize)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var categoryID *uint64
	if req.CategoryId > 0 {
		categoryID = &req.CategoryId
	}

	var videoStatus *domain.VideoStatus
	if req.Status != videopb.VideoStatus_VIDEO_STATUS_UNKNOWN {
		status := mapStatusFromProto(req.Status)
		videoStatus = &status
	}

	videos, total, err := s.videoService.ListVideos(ctx, &service.ListVideosInput{
		Page:       page,
		PageSize:   pageSize,
		CategoryID: categoryID,
		Status:     videoStatus,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list videos")
	}

	protoVideos := make([]*videopb.Video, len(videos))
	for i, video := range videos {
		protoVideos[i] = s.mapToProto(video)
	}

	// totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	return &videopb.ListVideosResponse{
		Videos: protoVideos,
		Pagination: &commonpb.Pagination{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Total:    total,
		},
	}, nil
}

func (s *VideoServer) UpdateVideo(ctx context.Context, req *videopb.UpdateVideoRequest) (*videopb.UpdateVideoResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	video, err := s.videoService.UpdateVideo(ctx, &service.UpdateVideoInput{
		ID:           req.Id,
		Title:        req.Title,
		Description:  req.Description,
		Status:       mapStatusFromProto(req.Status),
		ThumbnailURL: req.ThumbnailUrl,
		HLSManifest:  req.HlsManifest,
		Duration:     int(req.Duration),
	})

	if err != nil {
		if err == domain.ErrVideoNotFound {
			return nil, status.Error(codes.NotFound, "video not found")
		}
		return nil, status.Error(codes.Internal, "failed to update video")
	}

	return &videopb.UpdateVideoResponse{
		Video: s.mapToProto(video),
	}, nil
}

func (s *VideoServer) DeleteVideo(ctx context.Context, req *videopb.DeleteVideoRequest) (*commonpb.Response, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.videoService.DeleteVideo(ctx, req.Id)
	if err != nil {
		if err == domain.ErrVideoNotFound {
			return nil, status.Error(codes.NotFound, "video not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete video")
	}

	return &commonpb.Response{
		Success: true,
		Message: "video deleted successfully",
	}, nil
}

func (s *VideoServer) IncrementViewCount(ctx context.Context, req *videopb.IncrementViewCountRequest) (*commonpb.Response, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.videoService.IncrementViewCount(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to increment view count")
	}

	return &commonpb.Response{
		Success: true,
		Message: "view count incremented",
	}, nil
}

// Helper functions
func (s *VideoServer) mapToProto(video *domain.Video) *videopb.Video {
	return &videopb.Video{
		Id:           video.ID,
		Title:        video.Title,
		Description:  video.Description,
		UploadedBy:   video.UploadedBy,
		Duration:     int32(video.Duration),
		ThumbnailUrl: video.ThumbnailURL,
		HlsManifest:  video.HLSManifest,
		Status:       mapStatusToProto(video.Status),
		CategoryId:   video.CategoryID,
		ViewCount:    int32(video.ViewCount),
		CreatedAt:    timestamppb.New(video.CreatedAt),
		UpdatedAt:    timestamppb.New(video.UpdatedAt),
	}
}

func mapStatusToProto(status domain.VideoStatus) videopb.VideoStatus {
	switch status {
	case domain.StatusPending:
		return videopb.VideoStatus_VIDEO_STATUS_PENDING
	case domain.StatusTranscoding:
		return videopb.VideoStatus_VIDEO_STATUS_TRANSCODING
	case domain.StatusReady:
		return videopb.VideoStatus_VIDEO_STATUS_READY
	case domain.StatusFailed:
		return videopb.VideoStatus_VIDEO_STATUS_FAILED
	default:
		return videopb.VideoStatus_VIDEO_STATUS_UNKNOWN
	}
}

func mapStatusFromProto(status videopb.VideoStatus) domain.VideoStatus {
	switch status {
	case videopb.VideoStatus_VIDEO_STATUS_PENDING:
		return domain.StatusPending
	case videopb.VideoStatus_VIDEO_STATUS_TRANSCODING:
		return domain.StatusTranscoding
	case videopb.VideoStatus_VIDEO_STATUS_READY:
		return domain.StatusReady
	case videopb.VideoStatus_VIDEO_STATUS_FAILED:
		return domain.StatusFailed
	default:
		return domain.StatusPending
	}
}
