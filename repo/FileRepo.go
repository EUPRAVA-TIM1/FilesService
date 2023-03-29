package repo

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"log"
	"mime/multipart"
	"os"
	"path"
)

var UnrecognizedImageFormatError = errors.New("Image format is unsupported")

type FileRepo interface {
	SaveImage(image multipart.File) (string, error)
	SavePdf(pdf multipart.File) (string, error)
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

func (f fileRepo) SaveImage(newImage multipart.File) (string, error) {
	imageData, imageType, err := image.Decode(newImage)
	if err != nil {
		log.Println(err)
		return "", err
	}
	imageName := fmt.Sprintf("%s.%s", uuid.New().String(), imageType)

	file, err := os.Create(path.Join(f.filePath, imageName))
	defer file.Close()

	if err != nil {
		log.Println(err)
		return "", err
	}
	switch imageType {
	case "jpeg":
		options := jpeg.Options{Quality: 100}
		err = jpeg.Encode(file, imageData, &options)
	case "png":
		err = png.Encode(file, imageData)
	case "default":
		return "", UnrecognizedImageFormatError
	}
	if err != nil {
		log.Println(err)
		return "", err
	}

	return imageName, nil
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
