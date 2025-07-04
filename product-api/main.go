package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func getProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    p := Product{ID: id, Name: "placeholder", Price: 0}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(p)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/products/{id}", getProduct).Methods(http.MethodGet)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    log.Println("Product API listening on :8080")
    log.Fatal(srv.ListenAndServe())
}
