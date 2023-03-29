package service

import (
	"FileService/repo"
	"image"
	"mime/multipart"
	"os"
)

type FileService interface {
	SaveImage(image multipart.File) (string, error)
	SavePdf(newPdf multipart.File) (string, error)
	GetImage(name string) (image.Image, string, error)
	GetPdf(name string) (*os.File, error)
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
func (s fileService) GetImage(name string) (image.Image, string, error) {
	return s.repo.GetImage(name)
}

func (s fileService) GetPdf(name string) (*os.File, error) {
	return s.repo.GetPdf(name)
}
