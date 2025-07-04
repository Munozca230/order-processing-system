package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

type Customer struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Active bool `json:"active"`
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    c := Customer{ID: id, Name: "placeholder", Active: true}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(c)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/customers/{id}", getCustomer).Methods(http.MethodGet)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    log.Println("Customer API listening on :8080")
    log.Fatal(srv.ListenAndServe())
}
