package banco

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/lib/pq"
)

const (
	DB_ADDRESS = "localhost"
	DB_PORT    = 5432
	DB_DRIVER  = "postgres"
	DB_USER    = "postgres"
	//DB_PASSWORD = "postgres"
	DB_PASSWORD = "admin"
	DB_NAME     = "postgres"
)

//AbrirConexaoBD realiza a abertura de conexão com o DB
func AbrirConexaoBD() (*sql.DB, error) {
	postgresAddress := DB_ADDRESS
	if os.Getenv("POSTGRES_ADDRESS") != "" {
		postgresAddress = os.Getenv("POSTGRES_ADDRESS")
	}

	var postgresPort int64
	postgresPort = DB_PORT
	if os.Getenv("POSTGRES_PORT") != "" {
		postgresPort, _ = strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 0, 32)
	}

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresAddress, postgresPort, DB_USER, DB_PASSWORD, DB_NAME)
	fmt.Println(postgresAddress, postgresPort)
	fmt.Println(dbinfo)
	db, err := sql.Open(DB_DRIVER, dbinfo)
	if err != nil {
		return nil, fmt.Errorf("Não foi possível estabelecer conexão com o banco de dados. Detalhes: %s", err)
	}
	fmt.Println("A conexão com o banco de dados foi estabelecida")

	return db, nil
}

//AbrirTransacaoEPreparar abre transão no banco e prepara o insert
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
