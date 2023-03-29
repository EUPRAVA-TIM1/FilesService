package service

import (
	"FileService/repo"
	"mime/multipart"
)

type FileService interface {
	SaveImage(image multipart.File) (string, error)
	SavePdf(newPdf multipart.File) (string, error)
}

type fileService struct {
	repo repo.FileRepo
}

func NewFileService(fileRepo repo.FileRepo) FileService {
	return fileService{fileRepo}
}

func (s fileService) SaveImage(newImage multipart.File) (string, error) {
	return s.repo.SaveImage(newImage)
}

func (s fileService) SavePdf(newPdf multipart.File) (string, error) {
	return s.repo.SavePdf(newPdf)
}
