package http

import (
	"context"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"reflect"
	"time"
)

type StoppableListener interface {
	Shutdown(ctx context.Context) error
	ListenAndServe() error
}

type Server struct {
	StoppableListener
	martini.Router
	Logger *log.Logger
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := s.StoppableListener.Shutdown(ctx); err != nil {
		s.Logger.Printf("Server shutdown error: %s", err)
	}
}

func (s *Server) Run() error {
	if err := s.StoppableListener.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func NewServer(port string) *Server {
	m := martini.Classic()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: m,
	}
	return &Server{
		StoppableListener: s,
		Router:            m,
		Logger:            m.Injector.Get(reflect.TypeOf(&log.Logger{})).Interface().(*log.Logger),
	}
}
