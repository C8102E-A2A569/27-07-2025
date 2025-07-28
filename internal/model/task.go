package model

const (
	StatusCreated    = "created"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

type Task struct {
	ID          string
	Status      string
	Files       []FileItem
	FailedURLs  []string
	ArchivePath string
}

type FileItem struct {
	URL      string
	Filename string
	Data     []byte
}
