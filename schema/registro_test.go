package schema

import "testing"

const erro = "valor esperado %v, valor retornado %v"

func TestStringWithoutSpecialCharacter(t *testing.T) {
	valorEsperado := "12345"
	resultado := StringWithoutSpecialCharacter([]byte("1.2,3-4/5"))
	if resultado != valorEsperado {
		t.Errorf(erro, valorEsperado, resultado)
	}
}

func TestGetNullString(t *testing.T) {
	resultado := GetNullString("NULL")
	if resultado.Valid {
		t.Errorf(erro, false, true)
	}
}
