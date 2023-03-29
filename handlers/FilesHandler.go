package handlers

import (
	"FileService/repo"
	"FileService/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type FilesHandler interface {
	Init(r *mux.Router)
}

const (
	fileKey             = "file"
	badRequestMsg       = "Bad Request"
	contentType         = "Content-Type"
	imageJpeg           = "image/jpeg"
	imagePng            = "image/png"
	appPdf              = "application/pdf"
	internalSrvErrMsg   = "Internal server error"
	unsupportedMediaMsg = "Unsupported media type"
)

type filesHandler struct {
	maxImageSize int64
	maxPdfSize   int64
	fileService  service.FileService
}

func NewFilesHandler(s service.FileService, maxImgSize int64, maxPdfSize int64) FilesHandler {
	return filesHandler{fileService: s, maxPdfSize: maxPdfSize, maxImageSize: maxImgSize}
}

func (f filesHandler) Init(r *mux.Router) {
	r.StrictSlash(false)
	r.HandleFunc("/api/files/{file}", f.GetFile).Methods("GET")
	r.HandleFunc("/api/files", f.SaveFile).Methods("POST")
	http.Handle("/", r)
}

func (f filesHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["file"]
	var fileType string

	fileType = filepath.Ext(name)

	if fileType == ".pdf" {
		w.Header().Set(contentType, appPdf)
		file, err := f.fileService.GetPdf(name)
		defer file.Close()
		if err == repo.FileNotExistError {
			http.Error(w, fmt.Sprintf("Pdf %q doesn't exist", name), http.StatusNotFound)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, badRequestMsg, http.StatusBadRequest)
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			log.Println(err)
			http.Error(w, badRequestMsg, http.StatusBadRequest)
			return
		}
	} else {
		var image image.Image
		var err error
		var imageType string

		image, imageType, err = f.fileService.GetImage(name)
		if err == repo.FileNotExistError {
			http.Error(w, fmt.Sprintf("Image %q doesn't exist", name), http.StatusNotFound)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, badRequestMsg, http.StatusBadRequest)
			return
		}

		switch imageType {
		case "jpeg":
			options := jpeg.Options{Quality: 100}
			w.Header().Set(contentType, imageJpeg)
			err = jpeg.Encode(w, image, &options)
		case "png":
			w.Header().Set(contentType, imagePng)
			err = png.Encode(w, image)
		case "default":
			log.Println(err.Error())
			http.Error(w, internalSrvErrMsg, http.StatusInternalServerError)
			return
		}
	}
}

func (f filesHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(f.maxImageSize)
	if err != nil {
		log.Println(err)
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}
	_, cnt, _ := r.FormFile(fileKey)

	contentTyp := filepath.Ext(cnt.Filename)
	switch contentTyp {
	case ".pdf":
		f.savePdf(w, r)
	case ".jpg":
		f.saveImg(w, r)
	case ".png":
		f.saveImg(w, r)
	default:
		http.Error(w, unsupportedMediaMsg, http.StatusUnsupportedMediaType)
	}

}

func (f filesHandler) saveImg(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, f.maxImageSize)
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}

	image, _, err := r.FormFile(fileKey)
	if err != nil {
		log.Println(err)
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}
	name, err := f.fileService.SaveImage(image)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}
	jsonResponse(struct {
		Name string `json:"name"`
	}{Name: name}, w, http.StatusCreated)
}

func (f filesHandler) savePdf(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(f.maxPdfSize)
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}

	pdf, _, err := r.FormFile(fileKey)
	if err != nil {
		log.Println(err)
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}
	name, err := f.fileService.SavePdf(pdf)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}
	jsonResponse(struct {
		Name string `json:"name"`
	}{Name: name}, w, http.StatusCreated)
}

func jsonResponse(object interface{}, w http.ResponseWriter, status int) {
	w.Header().Set(contentType, "application/json")
	resp, err := json.Marshal(object)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if status != 0 {
		w.WriteHeader(status)
	}
	_, err = w.Write(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
