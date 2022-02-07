package main

import (
	"bytes"
	apphttp "github.com/d-ashesss/news-feed-bot/http"
	"github.com/d-ashesss/news-feed-bot/pkg/db/memory"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/go-martini/martini"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type AppTest struct {
	testHttpServer *httptest.Server
	httpServer     *apphttp.Server
	logger         *log.Logger
	logBuffer      *bytes.Buffer
	app            *App
}

func (t *AppTest) Request(method string, url string, body io.Reader, headers map[string]string) (int, []byte, error) {
	t.testHttpServer.Start()
	defer t.testHttpServer.Close()

	req, err := http.NewRequest(method, t.testHttpServer.URL+url, body)
	if err != nil {
		return 0, nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := t.testHttpServer.Client().Do(req)
	if err != nil {
		return 0, nil, err
	}
	got, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, got, nil
}

func NewAppTest() *AppTest {
	config := Config{}
	handler := martini.Classic()
	testServer := httptest.NewUnstartedServer(handler)
	buffer := bytes.NewBufferString("")
	logger := log.New(buffer, "", 0)
	handler.Map(logger)
	httpServer := &apphttp.Server{StoppableListener: testServer.Config, Router: handler, Logger: logger}

	userStore := memory.NewUserStore()
	userModel := model.NewUserModel(userStore)

	return &AppTest{
		testHttpServer: testServer,
		httpServer:     httpServer,
		logger:         logger,
		logBuffer:      buffer,
		app:            NewApp(config, httpServer, userModel, nil),
	}
}
