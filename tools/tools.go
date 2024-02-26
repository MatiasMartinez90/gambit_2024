package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Para transformar la fecha a como la maneja mysql
func FechaMySQL() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func EscapeString(t string) string {
	desc := strings.ReplaceAll(t, "'", "")
	desc = strings.ReplaceAll(desc, "\"", "")
	return desc
}

func ArmoSentencia(s string, fieldsName string, typeField string, ValueN int, ValueF float64, ValueS string) string {
	//confirma consistencia de los campos
	if (typeField == "S") && len(ValueS) == 0 ||
		(typeField == "F" && ValueF == 0) ||
		(typeField == "N" && ValueN == 0) {
		return s
	}

	if !strings.HasSuffix(s, "SET ") {
		s += ", "
	}

	switch typeField {
	case "S":
		s += fieldsName + " = '" + EscapeString(ValueS) + "'"
	case "N":
		s += fieldsName + " = " + strconv.Itoa(ValueN)
	case "F":
		s += fieldsName + " = " + strconv.FormatFloat(ValueF, 'e', -1, 64)
	}

	return s

}
