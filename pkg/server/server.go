package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nedokyrill/posts-service/pkg/logger"
)

type APIServer struct {
	httpServer *http.Server
}

func NewAPIServer(router *gin.Engine) *APIServer {
	return &APIServer{
		httpServer: &http.Server{
			Addr:         ":" + os.Getenv("API_PORT"),
			Handler:      router.Handler(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *APIServer) Start() {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatal(err)
	}

}

func (s *APIServer) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	switch ctx.Err() {
	case context.DeadlineExceeded:
		logger.Logger.Error("timeout shutting down server")
	case nil:
		logger.Logger.Info("shutdown completed before timeout.")
	default:
		logger.Logger.Error("shutdown ended with:", ctx.Err())
	}

	return nil
}
