package routers

import (
	"encoding/json"
	"fmt"

	"github.com/gambit/bd"
	"github.com/gambit/models"
)

func InsertAddress(body string, User string) (int, string) {
	fmt.Println("Inicializando funcion  router.InsertAddress")

	var a models.Address

	err := json.Unmarshal([]byte(body), &a)
	if err != nil {
		return 400, "Error en los datos recibidos de direcciones" + err.Error()
	}

	if a.AddAddress == "" {
		return 400, "Debe especificar La direccion "
	}

	if a.AddName == "" {
		return 400, "Debe especificar El name de la direccion "
	}

	if a.AddTitle == "" {
		return 400, "Debe especificar el titulo "
	}

	if a.AddCity == "" {
		return 400, "Debe especificar la ciudad "
	}

	if a.AddPhone == "" {
		return 400, "Debe especificar el telefono "
	}

	if a.AddPostalCode == "" {
		return 400, "Debe especificar el codigo postal "
	}

	err = bd.InsertAddress(a, User)
	if err != nil {
		return 400, "Ocurrio un error al intentar registrar la direccion para el ID de usuario" + User + " > " + err.Error()
	}

	return 200, "Insert Address OK "

}

func UpdateAddress(body string, User string, id int) (int, string) {
	fmt.Println("Inicializando funcion  router.UpdateAddress")

	var t models.Address

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos" + err.Error()
	}

	t.AddId = id
	var encontrado bool
	err, encontrado = bd.AddressExists(User, t.AddId)

	if !encontrado {
		if err != nil {
			return 400, "Error al intentar buscar direccion del usuario" + User + " > " + err.Error()
		}
		return 400, "No se encuentra el registro del ID usuario asociado" + User + " > " + err.Error()
	}

	err = bd.UpdateAddress(t)
	if err != nil {
		return 400, "Ocurrio un error al intentar actualizar la actualizacion de la direccion" + User + " > " + err.Error()
	}

	return 200, " Update OK"

}

func DeleteAddres(User string, id int) (int, string) {
	fmt.Println("Inicializando funcion  router.DeleteAddres")

	if id == 0 {
		return 400, "Debe especificar ID de la direccion a Eliminar "
	}

	err, encontrado := bd.AddressExists(User, id)
	if !encontrado {
		if err != nil {
			return 400, "Error al intentar buscar direccion del usuario" + User + " > " + err.Error()
		}
		return 400, "No se encuentra el registro del ID usuario asociado" + User + " > " + err.Error()
	}

	err = bd.DeleteAddress(id)
	if err != nil {
		return 400, "Ocurrio un error al intentar borrar los datos del usuario" + err.Error()
	}

	return 200, " Delete OK "

}

func SelectAddress(User string) (int, string) {
	fmt.Println("Inicializando funcion  router.SelectAddress")

	addr, err := bd.SelectAddress(User)

	if err != nil {
		return 400, "Ocurrio un error al intentar capturar Resultados de Address/s  del Usuario: > " + User + "' en productos > " + err.Error()
	}

	respJson, err := json.Marshal(addr)
	if err != nil {
		return 400, "Ocurrio un error al intentar convertir en JSON  la busqueda de productos/s > " + err.Error()
	}

	return 200, string(respJson)

}
