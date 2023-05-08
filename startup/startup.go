package startup

import (
	"FileService/config"
	"FileService/handlers"
	"FileService/repo"
	"FileService/service"
	"context"
	"fmt"
	muxHandlers "github.com/gorilla/handlers"
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
func (server Server) setup(path string, maxFileSize int64) handlers.FilesHandler {
	fileRepo := repo.NewFileRepo(path)
	fileService := service.NewFileService(fileRepo)
	filesHandler := handlers.NewFilesHandler(fileService, maxFileSize)
	return filesHandler
}

func (server Server) Start() {

	r := mux.NewRouter()

	corsHandler := muxHandlers.CORS(
		muxHandlers.AllowedOrigins([]string{"*"}),
		muxHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		muxHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	h := server.setup(server.config.FilesPath, server.config.MaxFileSize)
	h.Init(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: corsHandler(r),
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
	log.Printf("Max file size is %d bytes", server.config.MaxFileSize)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server %s", err)
	}
	log.Println("server gracefully stopped")

}
