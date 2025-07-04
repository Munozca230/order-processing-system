package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func getProduct(c echo.Context) error {
    id := c.Param("id")
    p := Product{ID: id, Name: "placeholder", Price: 0}
    return c.JSON(http.StatusOK, p)
}

func main() {
    e := echo.New()
    e.GET("/products/:id", getProduct)
    e.Logger.Fatal(e.Start(":8080"))
}
