package startup

import (
	"FileService/config"
	"FileService/handlers"
	"FileService/repo"
	"FileService/service"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}
func (server Server) setup(path string, maxImgSize int64, maxPdfSize int64) handlers.FilesHandler {
	fileRepo := repo.NewFileRepo(path)
	fileService := service.NewFileService(fileRepo)
	filesHandler := handlers.NewFilesHandler(fileService, maxImgSize, maxPdfSize, fileRepo)
	return filesHandler
}

func (server Server) Start() {

	r := mux.NewRouter()

	h := server.setup(server.config.FilesPath, server.config.MaxImageSize, server.config.MaxPdfSize)
	h.Init(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: r,
	}

	wait := time.Second * 15
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	log.Printf("Listening on port = %s\n", server.config.Port)
	log.Printf("Files will be saved to = %s\n", server.config.FilesPath)
	log.Printf("Max image size is %d bytes, max pdf size is %d bytes", server.config.MaxImageSize, server.config.MaxPdfSize)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server %s", err)
	}
	log.Println("server gracefully stopped")

}
