package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gorilla/mux"
)

func TestGetProduct(t *testing.T) {
    r := mux.NewRouter()
    r.HandleFunc("/products/{id}", getProduct).Methods(http.MethodGet)

    req := httptest.NewRequest(http.MethodGet, "/products/123", nil)
    resp := httptest.NewRecorder()

    r.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", resp.Code)
    }
}
