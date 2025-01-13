package integration_test

// This file mainly does integration testing for the product creation

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func beforeEachController() {}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World\n"))
}

func TestControllerCreate(t *testing.T) {
	beforeEachController()
	t.Run("ShouldCreateAValidProductWhenPassingAllValidParams", func(t *testing.T) {
		req := httptest.NewRequest("POST", "localhost:8080", nil)
		fmt.Println(req)
		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()

		fmt.Println(resp.StatusCode)
	})
}
