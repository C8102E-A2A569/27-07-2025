package model

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

func CreateArchive(task *Task) (string, error) {
	if len(task.Files) == 0 {
		return "", fmt.Errorf("no files to archive")
	}
	if err := os.MkdirAll("archives", 0755); err != nil {
		return "", err
	}
	archivePath := filepath.Join("archives", task.ID+".zip")

	file, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	defer writer.Close()

	for _, fileItem := range task.Files {
		zipFile, err := writer.Create(fileItem.Filename)
		if err != nil {
			return "", err
		}

		_, err = zipFile.Write(fileItem.Data)
	}
	return archivePath, nil
}
