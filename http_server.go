package main

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

var port string
var httpHost string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	httpHost = fmt.Sprintf(":%s", port)
}

type Server struct {
	*http.Server
	*martini.ClassicMartini
	logger *log.Logger
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		s.logger.Printf("Server shutdown error: %s", err)
	}
}

func (s *Server) Run() error {
	if martini.Env == martini.Dev {
		s.logger.Printf("Listening on port %s", port)
		s.logger.Printf("Open http://localhost:%s in the browser", port)
	}
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func NewHttpServer() Server {
	m := martini.Classic()
	server := &http.Server{Addr: httpHost}

	http.Handle("/", m)

	return Server{
		Server:         server,
		ClassicMartini: m,
		logger:         m.Injector.Get(reflect.TypeOf(&log.Logger{})).Interface().(*log.Logger),
	}
}
