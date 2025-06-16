package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/yuuhikaze/rgxr/handlers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/union", handlers.UnionHandler).Methods("POST")
	r.HandleFunc("/convert", handlers.ConvertHandler).Methods("POST")
	r.HandleFunc("/save-image", handlers.SaveImageHandler).Methods("POST")

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
