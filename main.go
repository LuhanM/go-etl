package main

import (
	// "bufio"
	// "bytes"
	// "database/sql"
	// "fmt"
	// "io"
	"log"
	"net/http"
	// "strconv"
	// "time"

	"github.com/LuhanM/go-etl/importador"
	"github.com/gorilla/mux"
	//"github.com/lib/pq"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/arquivo", importador.ImportarArquivo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
