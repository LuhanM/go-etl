package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	importador "github.com/LuhanM/go-etl/importador"
)

// ImportarArquivo é o Handle que deve ser utilziado para importar massa de dados
func ImportarArquivo(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	transacional := true
	if r.FormValue("transacional") != "" {
		transacional, _ = strconv.ParseBool(r.FormValue("transacional"))
	}
	defer file.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, file)
	reader := bufio.NewReader(buf)
	err = importador.IgnorarCabecalho(reader)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	start := time.Now()

	err = importador.ImportarRegistros(reader, transacional)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	fim := time.Now()

	w.Write([]byte(fmt.Sprint("\nHora do início     ", start)))
	w.Write([]byte(fmt.Sprint("\nHora da conclusão  ", fim)))
	w.Write([]byte(fmt.Sprint("\nTempo de execução  ", fim.Sub(start))))

	fmt.Println("Hora do início    ", start)
	fmt.Println("Hora da conclusão ", fim)
	fmt.Println("Tempo de execução", fim.Sub(start))
}
