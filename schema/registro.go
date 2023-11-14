package schema

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Nhanderu/brdoc"
	"github.com/lib/pq"
)

type Registro struct {
	Cpf                string
	CpfValido          bool
	Private            bool
	Incompleto         bool
	DataUltimaCompra   pq.NullTime
	TicketMedio        sql.NullFloat64
	TicketUltimaCompra sql.NullFloat64
	LojaMaisFrequente  sql.NullString
	LojaUltimaCompra   sql.NullString
}

// GetNullString Retorna a partir de uma string o struct com os dado e sua validação
func GetNullString(p string) sql.NullString {
	var result sql.NullString
	result.String = p
	result.Valid = strings.ToUpper(result.String) != "NULL"
	return result
}

// StringWithoutSpecialCharacter remover caracteres de pontuação do []byte e retorna em uma string
func StringWithoutSpecialCharacter(p []byte) string {
	content := p
	content = bytes.ReplaceAll(content, []byte("."), []byte(""))
	content = bytes.ReplaceAll(content, []byte(","), []byte(""))
	content = bytes.ReplaceAll(content, []byte("-"), []byte(""))
	content = bytes.ReplaceAll(content, []byte("/"), []byte(""))
	return string(content)
}

// ParseByteToFloat realiza parse de []byte para sql.NullFloat64
func ParseByteToFloat(p []byte) sql.NullFloat64 {
	var result sql.NullFloat64
	var err error
	content := bytes.ReplaceAll(p, []byte(","), []byte("."))
	result.Float64, err = strconv.ParseFloat(string(content), 2)
	result.Valid = (err == nil)
	return result
}

// FormatarRegistro recebe array de campos e converte para objeto estruturado
func FormatarRegistro(campos [][]byte) (Registro, error) {
	var reg Registro
	var err error
	reg.Cpf = StringWithoutSpecialCharacter(campos[0])
	reg.CpfValido = brdoc.IsCPF(reg.Cpf)
	reg.LojaMaisFrequente = GetNullString(StringWithoutSpecialCharacter(campos[6]))
	reg.LojaUltimaCompra = GetNullString(StringWithoutSpecialCharacter(campos[7]))
	reg.Private, err = strconv.ParseBool(string(campos[1]))
	if err != nil {
		return reg, fmt.Errorf("erro ao converter valor do campo 'Private'. Detalhes: %s", err)
	}

	reg.Incompleto, err = strconv.ParseBool(string(campos[2]))
	if err != nil {
		return reg, fmt.Errorf("erro ao converter valor do campo 'Incompleto'. Detalhes: %s", err)
	}

	reg.DataUltimaCompra.Time, err = time.Parse("2006-01-02", string(campos[3]))
	reg.DataUltimaCompra.Valid = (err == nil)

	reg.TicketMedio = ParseByteToFloat(campos[4])
	reg.TicketUltimaCompra = ParseByteToFloat(campos[5])

	return reg, nil
}
