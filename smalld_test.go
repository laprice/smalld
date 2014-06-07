package main

import "log"
import "testing"
import "net/http"
import "net/http/httptest"

func TestLocationHandlerResponseOK(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/location", nil)
	LocationHandler(response, request)
	log.Println("/location response:", response.Code)
	if response.Code != http.StatusNoContent {
		t.Fatalf("Bad Response")
	}

}

func TestLocationHandlerQueryParse(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/location?lat=44.09491559960329&lon=-123.0965916720434&acc=5&label=foo", nil)
	LocationHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Bad Response")
	}
	log.Println(response)
}
