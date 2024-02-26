package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Estructura que viaja en el token
type TokenJSON struct {
	Sub       string
	Event_Id  string
	Token_use string
	Scope     string
	Auth_time int
	Iss       string
	Exp       int
	Iat       int
	Client_id string
	Username  string
}

func ValidoToken(token string) (bool, error, string) {
	fmt.Println("Ingresando a la funcion ValidoToken")
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		fmt.Println("El token no es valido")
		return false, nil, "El token no es valido"
	}
	fmt.Println("Hello, Tute")
	//userInfo, err := base64.StdEncoding.Strict().DecodeString("eyJhdF9oYXNoIjoicWtOdW50bUxYNncxaEhfOThFTmxKZyIsInN1YiI6Ijk0MzgyNGY4LTYwZTEtNzBhNi0zODFlLWE4YjhhNjU0OWVlZCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb21cL3VzLWVhc3QtMV9EdmhKNjJJZHoiLCJjb2duaXRvOnVzZXJuYW1lIjoiOTQzODI0ZjgtNjBlMS03MGE2LTM4MWUtYThiOGE2NTQ5ZWVkIiwiYXVkIjoiNTNmYzQ0OTVvdWtqajA0dWRtdmpydmZiMWQiLCJldmVudF9pZCI6IjllMzVmNTM4LWQzNDMtNDI1YS1iOWRlLTJhMTI3MzFiMjZkZiIsInRva2VuX3VzZSI6ImlkIiwiYXV0aF90aW1lIjoxNjkyNTQ0OTcwLCJleHAiOjE2OTI2MzEzNzAsImlhdCI6MTY5MjU0NDk3MCwianRpIjoiMzJlNDNkYTktNmY2Ni00MmZmLTk0YmMtNTdhNzFmYzMwMGUzIiwiZW1haWwiOiJtYXRpYXMubWFydGluZXo5MCs4N0BnbWFpbC5jb20ifQ")
	userInfo3, err := base64.StdEncoding.Strict().DecodeString(parts[1])
	userInfo2 := string(userInfo3) + "}"
	userInfo := []byte(userInfo2)

	fmt.Println("Imprimo userInfo2:")
	fmt.Println(userInfo2)
	fmt.Println("Imprimo userInfo:")
	fmt.Println(userInfo)

	//userInfo, err := base64.StdEncoding.DecodeString(parts[1])
	fmt.Println(err)
	//if err != nil {
	//	fmt.Println("No se puede decodificar la parte del token : ", err.Error())
	//	return false, err, err.Error()
	//}

	var tkj TokenJSON
	err = json.Unmarshal(userInfo, &tkj)
	if err != nil {
		fmt.Println("No se puede decodificar en la estructura JSON", err.Error())
		return false, err, err.Error()
	}

	ahora := time.Now()
	//creo una variable para poder comparar las fechas
	tm := time.Unix(int64(tkj.Exp), 0)

	//Funcion para comparar fechas
	if tm.Before(ahora) {
		fmt.Println("Fecha de expiracion token" + tm.String())
		fmt.Println("Token expirado !")
		return false, err, "Token expierado !!"
	}

	return true, nil, string(tkj.Sub)

}
