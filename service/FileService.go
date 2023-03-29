package service

import (
	"FileService/repo"
	"mime/multipart"
	"os"
)

type FileService interface {
	SaveFile(newFile multipart.File, fileType string) (string, error)
	GetFile(name string) (*os.File, error)
}

type fileService struct {
	repo repo.FileRepo
}

func NewFileService(fileRepo repo.FileRepo) FileService {
	return fileService{fileRepo}
}

func (s fileService) SaveFile(newFile multipart.File, fileType string) (string, error) {
	return s.repo.SaveFile(newFile, fileType)
}

func (s fileService) GetFile(name string) (*os.File, error) {
	return s.repo.GetFile(name)
}
