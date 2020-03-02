package main

import (
	"log"
	"net/http"

	"github.com/LuhanM/go-etl/importador"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/arquivo", importador.ImportarArquivo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
