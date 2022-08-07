package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startBackend() {
	router := mux.NewRouter()

	router.HandleFunc("/login", CreateAuth).Methods("GET")

	fmt.Println("Server at http://localhost:3472")
	log.Fatal(http.ListenAndServe(":3472", router))
}

func CreateAuth(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("code")

	log.Println("Token: ", token)
	callbackTokenFromWeb(token)

	http.ServeFile(w, r, "./static/index.html")
}
