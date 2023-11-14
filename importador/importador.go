package importador

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"

	"github.com/LuhanM/go-etl/config"
	"github.com/LuhanM/go-etl/schema"
)

func ImportarRegistros(reader *bufio.Reader, usarTransacaoGeral bool) error {

	var err error
	db, err := config.AbrirConexaoBD()
	if err != nil {
		return fmt.Errorf("não foi possível realizar a conexão com o banco de dados. Detalhes: %s", err)
	}

	defer db.Close()

	var txn *sql.Tx
	var stmt *sql.Stmt

	if usarTransacaoGeral {
		txn, stmt, err = config.AbrirTransacaoEPreparar(db)
		if err != nil {
			return fmt.Errorf("não foi possível abrir a transação. Detalhes: %s", err)
		}
		defer stmt.Close()
		defer txn.Commit()
	}
	err = processarRegistros(reader, usarTransacaoGeral, db, txn, stmt)
	if err != nil {
		return err
	}
	if usarTransacaoGeral {
		_, err := stmt.Exec()
		if err != nil {
			return fmt.Errorf("não foi possível finalizar a peraparação para inclusão na base de dados. Detalhes: %s", err)
		}
	}
	return nil
}

func IgnorarCabecalho(reader *bufio.Reader) error {
	_, _, err := reader.ReadLine()
	if err == io.EOF {
		return fmt.Errorf("arquivo está vazio")
	}
	return nil
}

func processarRegistros(reader *bufio.Reader, usarTransacaoGeral bool, db *sql.DB, txn *sql.Tx, stmt *sql.Stmt) error {
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

		reg, err := schema.FormatarRegistro(campos)
		if err != nil {
			fmt.Printf("Erro ao formatar registro. Detalhes: %s \n registro: %q", err, campos)
			continue
		}

		if usarTransacaoGeral {
			_, err = stmt.Exec(reg.Cpf, reg.CpfValido, reg.DataUltimaCompra, reg.Incompleto, reg.LojaMaisFrequente, reg.LojaUltimaCompra, reg.Private, reg.TicketMedio, reg.TicketUltimaCompra)
			if err != nil {
				return fmt.Errorf("erro ao inserir registro no bulk da transação. Detalhes %s", err)
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
