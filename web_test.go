package main

import "testing"

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
