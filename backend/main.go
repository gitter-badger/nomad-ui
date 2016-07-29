package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/ws", ws)

	log.Fatal(http.ListenAndServe(":6464", router))
}
