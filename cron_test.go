package main

import (
	"net/http"
	"testing"
)

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
