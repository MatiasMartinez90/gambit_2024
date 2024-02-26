package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gambit/bd"
	"github.com/gambit/models"
)

func InsertCategory(body string, User string) (int, string) {
	var t models.Category

	err := json.Unmarshal([]byte(body), &t)

	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	if len(t.CategName) == 0 {
		return 400, "Debe especificar el Nombre (Title) de la categoria"
	}

	if len(t.CategPath) == 0 {
		return 400, "Debe especificar el Path (Ruta) de la categoria"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		return 400, msg
	}

	result, err2 := bd.InsertCategory(t)
	if err2 != nil {
		return 400, "Ocurrio un error al intentar realizar el registro de la categoria " + t.CategName + " > " + err2.Error()
	}

	return 200, "{ CategID: " + strconv.Itoa(int(result)) + "}"
}

func UpdateCategory(body string, User string, id int) (int, string) {
	var t models.Category

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos" + err.Error()
	}

	if len(t.CategName) == 0 && len(t.CategPath) == 0 {
		return 400, "Debe especificar CategName y CategPath para actualizar "
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		return 400, msg
	}

	t.CategID = id
	err2 := bd.UpdateCategory(t)
	if err2 != nil {
		return 400, "Ocurrio un error al intentar realizar el UPDATE de la categoria " + strconv.Itoa(id) + " > " + err.Error()
	}

	return 200, "Update OK"
}

func DeleteCategory(body string, User string, id int) (int, string) {
	if id == 0 {
		return 400, "Debe especificar ID de la Categoria a borrar"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		return 400, msg
	}

	err := bd.DeleteCategory(id)
	if err != nil {
		return 400, "Ocurrio un error al intentar realizar el DELETE de la categoria " + strconv.Itoa(id) + " > " + err.Error()
	}

	return 200, "Delete OK"
}

func SelectCategories(body string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Inicializando funcion  router.SelectCategories")

	var err error
	var CategId int
	var Slug string

	if len(request.QueryStringParameters["categId"]) > 0 {
		CategId, err = strconv.Atoi(request.QueryStringParameters["categId"])
		if err != nil {
			return 500, "Ocurrio un error al intentar convertir el valor" + request.QueryStringParameters["categId"]
		}

	} else {
		if len(request.QueryStringParameters["slug"]) > 0 {
			Slug = request.QueryStringParameters["slug"]
		}
	}

	lista, err2 := bd.SelectCategories(CategId, Slug)
	if err2 != nil {
		return 400, "Ocurrio un error al intentar capturar Categoria/s > " + err2.Error()
	}

	Categ, err3 := json.Marshal(lista)
	if err3 != nil {
		return 400, "Ocurrio un error al intentar convertir en JSON  Categoria/s > " + err3.Error()
	}

	return 200, string(Categ)

}
