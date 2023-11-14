package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	importador "github.com/LuhanM/go-etl/importador"
	"github.com/gin-gonic/gin"
)

func ImportFile(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.Writer.WriteHeader(400)
		ctx.Writer.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	transacional := true
	if ctx.Request.FormValue("transacional") != "" {
		transacional, _ = strconv.ParseBool(ctx.Request.FormValue("transacional"))
	}
	defer file.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, file)
	reader := bufio.NewReader(buf)
	err = importador.IgnorarCabecalho(reader)
	if err != nil {
		ctx.Writer.WriteHeader(400)
		ctx.Writer.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	start := time.Now()

	err = importador.ImportarRegistros(reader, transacional)
	if err != nil {
		ctx.Writer.WriteHeader(400)
		ctx.Writer.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	fim := time.Now()

	ctx.Writer.Write([]byte(fmt.Sprint("\nHora do início     ", start)))
	ctx.Writer.Write([]byte(fmt.Sprint("\nHora da conclusão  ", fim)))
	ctx.Writer.Write([]byte(fmt.Sprint("\nTempo de execução  ", fim.Sub(start))))

	fmt.Println("Hora do início    ", start)
	fmt.Println("Hora da conclusão ", fim)
	fmt.Println("Tempo de execução", fim.Sub(start))
}
