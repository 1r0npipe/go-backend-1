package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadHandler_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/?ext=jpg", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := &UploadHandler{
		UploadDir: fileSystem,
	}
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// expected := string(`[{"filename":"test.jpg","sizeByte":51358}]`)
	// if rr.Result().StatusCode != ht {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
