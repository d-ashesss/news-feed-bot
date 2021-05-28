package main

import (
	"bytes"
	"context"
	"log"
	"testing"
)

type StateSever struct {
	Running bool
}

func (s *StateSever) ListenAndServe() error {
	s.Running = true
	return nil
}

func (s *StateSever) Shutdown(_ context.Context) error {
	s.Running = false
	return nil
}

type ServerTest struct {
	state     *StateSever
	server    *Server
	logger    *log.Logger
	logBuffer *bytes.Buffer
}

func NewServerTest(initialState bool) *ServerTest {
	state := &StateSever{Running: initialState}
	buffer := bytes.NewBufferString("")
	logger := log.New(buffer, "", 0)
	server := &Server{HttpServer: state, Logger: logger}
	return &ServerTest{
		state:     state,
		server:    server,
		logger:    logger,
		logBuffer: buffer,
	}
}

func TestServer_Run(t *testing.T) {
	test := NewServerTest(false)
	_ = test.server.Run()
	if true != test.state.Running {
		t.Errorf("Test server hasn't been started")
	}
}

func TestServer_Shutdown(t *testing.T) {
	test := NewServerTest(true)
	test.server.Shutdown()
	if false != test.state.Running {
		t.Errorf("Test server hasn't been stopped")
	}
}
