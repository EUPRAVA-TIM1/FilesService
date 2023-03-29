package repo

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"

	"log"
	"mime/multipart"
	"os"
	"path"
)

var UnrecognizedFormatError = errors.New("File format is unsupported")
var FileNotExistError = errors.New("File doesnt exists")

type FileRepo interface {
	SaveFile(image multipart.File, fileType string) (string, error)
	GetFile(name string) (*os.File, error)
}
type fileRepo struct {
	filePath string
}

func NewFileRepo(path string) FileRepo {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println(err)
		err = os.MkdirAll(path, 0770)
		if err != nil {
			log.Println(err)
		}
	} else if err != nil {
		log.Println(err)
	}
	return fileRepo{path}
}

func (f fileRepo) SavePdf(pdf multipart.File) (string, error) {
	pdfName := fmt.Sprintf("%s.%s", uuid.New().String(), "pdf")
	file, err := os.Create(path.Join(f.filePath, pdfName))
	defer file.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = io.Copy(file, pdf)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return pdfName, nil

}

func (f fileRepo) SaveFile(newFile multipart.File, fileType string) (string, error) {
	if fileType != "jpeg" && fileType != "png" && fileType != "pdf" {
		return "", UnrecognizedFormatError
	}

	fileName := fmt.Sprintf("%s.%s", uuid.New().String(), fileType)
	file, err := os.Create(path.Join(f.filePath, fileName))
	defer file.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = io.Copy(file, newFile)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fileName, nil
}

func (f fileRepo) GetFile(name string) (*os.File, error) {
	_, err := os.Stat(path.Join(f.filePath, name))
	if os.IsNotExist(err) {
		return nil, FileNotExistError
	}
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path.Join(f.filePath, name))
	if err != nil {
		return nil, err
	}

	return file, nil
}
