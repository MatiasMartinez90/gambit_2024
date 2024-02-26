package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gambit/bd"
	"github.com/gambit/models"
)

func UpdateUser(body string, User string) (int, string) {
	fmt.Println("Inicializando funcion  router.UpdateUser")

	var t models.User
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos" + err.Error()
	}

	if len(t.UserFirstName) == 0 && len(t.UserLastName) == 0 {
		return 400, "Debe especificar el nombre (FirstName) o (LastName) del Usuario " + err.Error()
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		return 400, msg
	}

	_, encontrado := bd.UserExists(User)
	if !encontrado {
		return 400, "No existe un usuario con ese UUID " + User + "'"
	}

	err = bd.UpdateUser(t, User)
	if err != nil {
		return 400, "Ocurrio un error al intentar realizar la actualización del usuario " + User + " > " + err.Error()
	}

	return 200, " UpdatUser OK "
}

func SelectUser(body string, User string) (int, string) {
	fmt.Println("Inicializando funcion  router.SelectUser")

	_, encontrado := bd.UserExists(User)
	if !encontrado {
		return 400, "No existe un usuario con ese UUID " + User + "'"
	}

	row, err := bd.SelectUser(User)
	fmt.Println(row)
	if err != nil {
		return 400, "Ocurrio un error al intentar realizar el select del usuario " + User + " > " + err.Error()
	}

	respJson, err := json.Marshal(row)
	if err != nil {
		return 500, "Error al formatear los datos del usuario como JSON"
	}

	return 200, string(respJson)
}

func SelectUsers(body string, User string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Inicializando funcion  router.SelectUsers")

	var Page int
	if len(request.QueryStringParameters["page"]) == 0 {
		Page = 1
	} else {
		Page, _ = strconv.Atoi(request.QueryStringParameters["page"])
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		return 400, msg
	}

	user, err := bd.SelectUsers(Page)
	if err != nil {
		return 400, "Ocurrio un error al intentar obtener la lista de usuarios" + err.Error()
	}

	respJson, err := json.Marshal(user)
	if err != nil {
		return 500, "Ocurrio un error al intentar formatera en JSON los datos de los usarios > " + err.Error()
	}

	return 200, string(respJson)

}
