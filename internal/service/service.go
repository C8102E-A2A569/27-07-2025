package service

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"zip-archive/internal/config"
	"zip-archive/internal/model"
	proto "zip-archive/protos"
)

type Service struct {
	proto.UnsafeArchiveZipServiceServer
	tasks      map[string]*model.Task
	mu         sync.RWMutex
	downloader *model.Downloader
	config     *config.Config
}

func (s *Service) AddFiles(ctx context.Context, request *proto.AddFilesRequest) (*proto.AddFilesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetTaskStatus(ctx context.Context, request *proto.GetTaskStatusRequest) (*proto.GetTaskStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DownloadArchive(ctx context.Context, request *proto.DownloadArchiveRequest) (*proto.DownloadArchiveResponse, error) {
	//TODO implement me
	panic("implement me")
}

func New(cfg *config.Config) *Service {
	return &Service{
		tasks:      make(map[string]*model.Task),
		downloader: model.New(&cfg.Files),
		config:     cfg,
	}
}

func (s *Service) CreateTask(ctx context.Context, req *proto.CreateTaskRequest) (*proto.CreateTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.tasks) >= s.config.Limits.MaxTasks {
		return nil, status.Error(codes.ResourceExhausted, "server is busy")
	}

	taskID := uuid.New().String()
	task := &model.Task{
		ID:     taskID,
		Status: model.StatusCreated,
		Files:  make([]model.FileItem, 0),
	}

	s.tasks[taskID] = task

	return &proto.CreateTaskResponse{
		TaskId: taskID,
		Status: model.StatusCreated,
	}, nil
}

//func (s *Service) AddFiles(ctx context.Context)
