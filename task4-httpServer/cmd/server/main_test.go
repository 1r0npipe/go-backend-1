package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadHandler_ServeHTTP(t *testing.T) {
	testsTab := []struct {
		wantRequest string
		wantParam   string
		gotCode     int
		gotOutput   string
	}{
		{
			wantRequest: "GET",
			wantParam:   "/?ext=jpg",
			gotCode:     200,
			gotOutput:   `[{"filename":"test.jpg","sizeByte":51358}]`,
		},
		{
			wantRequest: "GET",
			wantParam:   "/?ext=txt",
			gotCode:     200,
			gotOutput:   `[{"filename":"test.txt","sizeByte":18},{"filename":"test1.txt","sizeByte":5}]`,
		},
		{
			wantRequest: "GET",
			wantParam:   "/ls",
			gotCode:     200,
			gotOutput:   `[{"filename":"test.jpg","sizeByte":51358},{"filename":"test.txt","sizeByte":18},{"filename":"test1.txt","sizeByte":5}]`,
		},
		{
			wantRequest: "GET",
			wantParam:   "/?ext=test",
			gotCode:     200,
			gotOutput:   `null`,
		},
	}
	for _, tt := range testsTab {
		req, err := http.NewRequest(tt.wantRequest, tt.wantParam, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := &UploadHandler{
			UploadDir: fileSystem,
		}
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != tt.gotCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		got := rr.Body.String()[:len(rr.Body.String())-1]
		if got != tt.gotOutput {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), tt.gotOutput)
		}
	}
}
