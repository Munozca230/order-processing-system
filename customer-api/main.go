package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

type Customer struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Active bool `json:"active"`
}

func getCustomer(c echo.Context) error {
    id := c.Param("id")
    customer := Customer{ID: id, Name: "placeholder", Active: true}
    return c.JSON(http.StatusOK, customer)
}

func main() {
    e := echo.New()
    e.GET("/customers/:id", getCustomer)
    e.Logger.Fatal(e.Start(":8080"))
}
