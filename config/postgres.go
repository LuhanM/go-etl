package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/lib/pq"
)

const (
	DB_ADDRESS  = "localhost"
	DB_PORT     = 5432
	DB_DRIVER   = "postgres"
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

// AbrirConexaoBD realiza a abertura de conexão com o DB
func AbrirConexaoBD() (*sql.DB, error) {
	dbAddress := os.Getenv("POSTGRES_ADDRESS")
	if dbAddress == "" {
		dbAddress = DB_ADDRESS
	}

	var dbPort int64
	dbPort = DB_PORT
	if os.Getenv("POSTGRES_PORT") != "" {
		dbPort, _ = strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 0, 32)
	}

	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		dbUser = DB_USER
	}

	dbPass := os.Getenv("POSTGRES_PASS")
	if dbPass == "" {
		dbPass = DB_PASSWORD
	}

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbAddress, dbPort, dbUser, dbPass, DB_NAME)
	fmt.Println(dbAddress, dbPort)
	fmt.Println(dbinfo)
	db, err := sql.Open(DB_DRIVER, dbinfo)
	if err != nil {
		return nil, fmt.Errorf("Não foi possível estabelecer conexão com o banco de dados. Detalhes: %s", err)
	}
	fmt.Println("A conexão com o banco de dados foi estabelecida")

	return db, nil
}

// AbrirTransacaoEPreparar abre transão no banco e prepara o insert
func AbrirTransacaoEPreparar(db *sql.DB) (*sql.Tx, *sql.Stmt, error) {
	txn, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := txn.Prepare(pq.CopyIn("estatistica", "cpf", "cpfvalido", "dataultimacompra", "incompleto", "lojamaisfrequente", "lojaultimacompra", "private", "ticketmedio", "ticketultimacompra"))
	if err != nil {
		return nil, nil, err
	}
	return txn, stmt, nil
}
