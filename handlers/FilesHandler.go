package handlers

import (
	"FileService/repo"
	"FileService/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
	maxFileSize int64
	fileService service.FileService
}

func NewFilesHandler(s service.FileService, maxFileSize int64) FilesHandler {
	return filesHandler{fileService: s, maxFileSize: maxFileSize}
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
	fileType := filepath.Ext(name)

	switch fileType {
	case ".pdf":
		w.Header().Set(contentType, appPdf)
	case ".jpeg":
		w.Header().Set(contentType, imageJpeg)
	case ".png":
		w.Header().Set(contentType, imagePng)
	default:
		http.Error(w, internalSrvErrMsg, http.StatusInternalServerError)
		return
	}

	file, err := f.fileService.GetFile(name)
	defer file.Close()
	if err == repo.FileNotExistError {
		http.Error(w, fmt.Sprintf("File %q doesn't exist", name), http.StatusNotFound)
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
		http.Error(w, internalSrvErrMsg, http.StatusInternalServerError)
		return
	}
}

func (f filesHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(f.maxFileSize)
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
	file, cnt, err := r.FormFile(fileKey)
	if err != nil {
		log.Println(err)
		http.Error(w, badRequestMsg, http.StatusBadRequest)
		return
	}

	contentTyp := filepath.Ext(cnt.Filename)
	switch contentTyp {
	case ".pdf":
		name, err := f.fileService.SaveFile(file, "pdf")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, badRequestMsg, http.StatusBadRequest)
			return
		}
		jsonResponse(struct {
			Name string `json:"name"`
		}{Name: name}, w, http.StatusCreated)
	case ".jpg":
		name, err := f.fileService.SaveFile(file, "jpeg")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, badRequestMsg, http.StatusBadRequest)
			return
		}
		jsonResponse(struct {
			Name string `json:"name"`
		}{Name: name}, w, http.StatusCreated)
	case ".png":
		name, err := f.fileService.SaveFile(file, "png")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, badRequestMsg, http.StatusBadRequest)
			return
		}
		jsonResponse(struct {
			Name string `json:"name"`
		}{Name: name}, w, http.StatusCreated)
	default:
		http.Error(w, unsupportedMediaMsg, http.StatusUnsupportedMediaType)
	}

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
