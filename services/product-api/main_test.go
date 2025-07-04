package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
)

func TestGetProduct(t *testing.T) {
    e := echo.New()

    req := httptest.NewRequest(http.MethodGet, "/products/123", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("123")

    if err := getProduct(c); err != nil {
        t.Fatalf("handler returned error: %v", err)
    }

    if rec.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", rec.Code)
    }
}
