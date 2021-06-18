package http

import (
	"context"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

var (
	port, httpHost string
)

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	httpHost = fmt.Sprintf(":%s", port)
}

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
	if martini.Env == martini.Dev {
		s.Logger.Printf("Listening on port %s", port)
		s.Logger.Printf("Open http://localhost:%s in the browser", port)
	}
	if err := s.StoppableListener.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func NewServer() *Server {
	m := martini.Classic()
	s := &http.Server{
		Addr:    httpHost,
		Handler: m,
	}
	return &Server{
		StoppableListener: s,
		Router:            m,
		Logger:            m.Injector.Get(reflect.TypeOf(&log.Logger{})).Interface().(*log.Logger),
	}
}
