package main

import (
	"bytes"
	"github.com/go-martini/martini"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AppTest struct {
	testHttpServer *httptest.Server
	httpServer     *Server
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
	handler := martini.Classic()
	testServer := httptest.NewUnstartedServer(handler)
	buffer := bytes.NewBufferString("")
	logger := log.New(buffer, "", 0)
	handler.Map(logger)
	httpServer := &Server{ClassicMartini: handler, HttpServer: testServer.Config, Logger: logger}
	return &AppTest{
		testHttpServer: testServer,
		httpServer:     httpServer,
		logger:         logger,
		logBuffer:      buffer,
		app:            NewApp(httpServer),
	}
}

func TestApp_authCron(t *testing.T) {
	testMethod := "GET"
	testUrl := "/cron-endpoint"
	tests := []struct {
		name       string
		headers    map[string]string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Unauthorized",
			headers:    map[string]string{},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "",
		},
		{
			name: "Authorized",
			headers: map[string]string{
				"X-Appengine-Cron": "true",
			},
			wantStatus: http.StatusOK,
			wantBody:   "ok",
		},
	}

	for _, testData := range tests {
		t.Run(testData.name, func(t *testing.T) {
			test := NewAppTest()
			test.httpServer.AddRoute(testMethod, testUrl, test.app.authCron, func() string {
				return testData.wantBody
			})

			gotStatus, gotBody, err := test.Request(testMethod, testUrl, nil, testData.headers)
			if err != nil {
				t.Fatal(err)
			}

			if gotStatus != testData.wantStatus {
				t.Errorf("got response %d, want %d", gotStatus, testData.wantStatus)
			}
			if string(gotBody) != testData.wantBody {
				t.Errorf("got response %q, want %q", gotBody, testData.wantBody)
			}
		})
	}
}

func TestApp_handleIndex(t *testing.T) {
	test := NewAppTest()

	gotStatus, gotBody, err := test.Request("GET", test.testHttpServer.URL+"/", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	wantStatus := 200
	if gotStatus != wantStatus {
		t.Errorf("got response %d, want %d", gotStatus, wantStatus)
	}
	wantBody := "Hello, World!"
	if string(gotBody) != wantBody {
		t.Errorf("got response %q, want %q", gotBody, wantBody)
	}
}
