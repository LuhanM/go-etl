package importador

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/LuhanM/go-etl/banco"
	"github.com/LuhanM/go-etl/util"
)

//ImportarArquivo é o Handle que deve ser utilziado para importar massa de dados
func ImportarArquivo(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	transacional = true
	if r.FormValue("transacional") != "" {
		transacional, _ = strconv.ParseBool(r.FormValue("transacional"))
	}
	defer file.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, file)
	reader := bufio.NewReader(buf)
	err = ignorarCabecalho(reader)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprint("Não foi possível importar os dados. Detalhe: ", err)))
	}

	start := time.Now()

	err = importarRegistros(reader)
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

func importarRegistros(reader *bufio.Reader) error {

	var err error
	db, err := banco.AbrirConexaoBD()
	if err != nil {
		return fmt.Errorf("Não foi possível realizar a conexão com o banco de dados. Detalhes: %s", err)
	}

	defer db.Close()

	var txn *sql.Tx
	var stmt *sql.Stmt

	if deveUsarTransacaoGeral() {
		txn, stmt, err = banco.AbrirTransacaoEPreparar(db)
		if err != nil {
			return fmt.Errorf("Não foi possível abrir a transação. Detalhes: %s", err)
		}
		defer stmt.Close()
		defer txn.Commit()
	}
	err = processarRegistros(reader, db, txn, stmt)
	if err != nil {
		return err
	}
	if deveUsarTransacaoGeral() {
		_, err := stmt.Exec()
		if err != nil {
			return fmt.Errorf("Não foi possível finalizar a peraparação para inclusão na base de dados. Detalhes: %s", err)
		}
	}
	return nil
}

var transacional bool

func deveUsarTransacaoGeral() bool {
	return transacional
}
func ignorarCabecalho(reader *bufio.Reader) error {
	_, _, err := reader.ReadLine()
	if err == io.EOF {
		return fmt.Errorf("Arquivo está vazio")
	}
	return nil
}

func processarRegistros(reader *bufio.Reader, db *sql.DB, txn *sql.Tx, stmt *sql.Stmt) error {
	var contador int
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		campos := bytes.Fields(line)
		if len(campos) != 8 {
			continue
		}

		reg, err := util.FormatarRegistro(campos)
		if err != nil {
			fmt.Printf("Erro ao formatar registro. Detalhes: %s \n registro: %q", err, campos)
			continue
		}

		if deveUsarTransacaoGeral() {
			_, err = stmt.Exec(reg.Cpf, reg.CpfValido, reg.DataUltimaCompra, reg.Incompleto, reg.LojaMaisFrequente, reg.LojaUltimaCompra, reg.Private, reg.TicketMedio, reg.TicketUltimaCompra)
			if err != nil {
				return fmt.Errorf("Erro ao inserir registro no bulk da transação. Detalhes %s", err)
			}
		} else {
			_ = db.QueryRow("insert into estatistica (cpf, cpfvalido, dataultimacompra, incompleto, lojamaisfrequente, lojaultimacompra, private, ticketmedio, ticketultimacompra) "+
				"values ($1, $2, $3, $4, $5, $6, $7, $8, $9);",
				reg.Cpf, reg.CpfValido, reg.DataUltimaCompra, reg.Incompleto, reg.LojaMaisFrequente, reg.LojaUltimaCompra, reg.Private, reg.TicketMedio, reg.TicketUltimaCompra).Scan()
		}

		contador++
		if (contador % 1000) == 0 {
			fmt.Println(contador, "registros encontrados")
		}
	}
	fmt.Println(contador, "registros encontrados")
	return nil
}
