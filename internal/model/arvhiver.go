package model

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

func CreateArchive(task *Task) (string, error) {
	if len(task.Files) == 0 {
		return "", fmt.Errorf("no files")
	}
	os.MkdirAll("archives", 0755)
	archivePath := filepath.Join("archives", task.ID+".zip")

	file, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	wr := zip.NewWriter(file)
	defer wr.Close()

	for _, f := range task.Files {
		zipFile, err := wr.Create(f.Filename)
		if err != nil {
			return "", err
		}
		zipFile.Write(f.Data)
	}
	return archivePath, nil
}
