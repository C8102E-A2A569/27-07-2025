package service

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
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

func (s *Service) AddFiles(ctx context.Context, req *proto.AddFilesRequest) (*proto.AddFilesResponse, error) {
	s.mu.Lock()
	task, exists := s.tasks[req.TaskId]
	if !exists {
		s.mu.Unlock()
		return nil, status.Error(codes.NotFound, "task not found")
	}

	if task.Status != model.StatusCreated {
		s.mu.Unlock()
		return nil, status.Error(codes.FailedPrecondition, "task isn't in created status")
	}

	if len(task.Files)+len(req.Urls) > s.config.Limits.MaxFilesPerTask {
		s.mu.Unlock()
		return nil, status.Error(codes.InvalidArgument, "too many files")
	}

	task.Status = model.StatusProcessing
	s.mu.Unlock()
	var failedURLs []string

	for _, url := range req.Urls {
		data, filename, err := s.downloader.Download(ctx, url)
		if err != nil {
			failedURLs = append(failedURLs, url)
			continue
		}

		s.mu.Lock()
		task.Files = append(task.Files, model.FileItem{
			URL:      url,
			Filename: filename,
			Data:     data,
		})
		s.mu.Unlock()
	}

	s.mu.Lock()
	task.FailedURLs = append(task.FailedURLs, failedURLs...)

	if len(task.Files) >= s.config.Limits.MaxFilesPerTask {
		go s.processTask(task.ID)
	} else {
		task.Status = model.StatusCreated
	}
	s.mu.Unlock()

	return &proto.AddFilesResponse{
		Status:     "ok",
		FailedUrls: failedURLs,
	}, nil
}

func (s *Service) processTask(taskID string) {
	s.mu.Lock()
	task := s.tasks[taskID]
	s.mu.Unlock()

	archivePath, err := model.CreateArchive(task)

	s.mu.Lock()
	if err != nil {
		task.Status = model.StatusFailed
	} else {
		task.Status = model.StatusCompleted
		task.ArchivePath = archivePath
	}
	s.mu.Unlock()
}

func (s *Service) GetTaskStatus(ctx context.Context, req *proto.GetTaskStatusRequest) (*proto.GetTaskStatusResponse, error) {
	s.mu.RLock()
	task, exists := s.tasks[req.TaskId]
	if !exists {
		s.mu.RUnlock()
		return nil, status.Error(codes.NotFound, "task not found")
	}

	response := &proto.GetTaskStatusResponse{
		TaskId:     task.ID,
		Status:     task.Status,
		FailedUrls: task.FailedURLs,
	}
	s.mu.RUnlock()

	return response, nil
}

func (s *Service) DownloadArchive(ctx context.Context, req *proto.DownloadArchiveRequest) (*proto.DownloadArchiveResponse, error) {
	s.mu.RLock()
	task, exists := s.tasks[req.TaskId]
	if !exists {
		s.mu.RUnlock()
		return nil, status.Error(codes.NotFound, "task not found")
	}

	if task.Status != model.StatusCompleted || task.ArchivePath == "" {
		s.mu.RUnlock()
		return nil, status.Error(codes.FailedPrecondition, "archive isn't ready")
	}
	s.mu.RUnlock()

	data, err := os.ReadFile(task.ArchivePath)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to read archive")
	}

	return &proto.DownloadArchiveResponse{
		ArchiveData: data,
		Filename:    task.ID + ".zip",
	}, nil
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

	activeTaskCnt := 0
	for _, task := range s.tasks {
		if task.Status == model.StatusCreated || task.Status == model.StatusProcessing {
			activeTaskCnt++
		}
	}
	if activeTaskCnt >= s.config.Limits.MaxTasks {
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
