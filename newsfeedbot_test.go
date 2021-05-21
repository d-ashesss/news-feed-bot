package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCronAuth(t *testing.T) {
	testMethod := "GET"
	testUrl := "/cron-endpoint"
	tests := []struct {
		headers    map[string]string
		wantStatus int
		wantBody   string
	}{
		{
			headers:    map[string]string{},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "",
		},
		{
			headers: map[string]string{
				"X-Appengine-Cron": "true",
			},
			wantStatus: http.StatusOK,
			wantBody:   "ok",
		},
	}

	for _, test := range tests {
		server := NewHttpServer()
		server.AddRoute(testMethod, testUrl, cronAuth, func() string {
			return "ok"
		})
		logger := log.New(&bytes.Buffer{}, "", 0)
		server.Map(logger)

		req := httptest.NewRequest(testMethod, testUrl, nil)
		for key, value := range test.headers {
			req.Header.Set(key, value)
		}
		rr := httptest.NewRecorder()

		server.ServeHTTP(rr, req)

		if got := rr.Result().StatusCode; got != test.wantStatus {
			t.Errorf("got status code %d, want %d", got, test.wantStatus)
		}
		if got := rr.Body.String(); got != test.wantBody {
			t.Errorf("got response body %q, want %q", got, test.wantBody)
		}
	}
}
