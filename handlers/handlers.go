package handlers

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gambit/auth"
	"github.com/gambit/routers"
)

func Manejadores(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Voy a procesar"+path, " > "+method)
	fmt.Println("Voy a procesar"+path[1:5], " > "+method)

	id := request.PathParameters["id"]
	idn, _ := strconv.Atoi(id)

	isOk, statusCode, user := validoAuthorization(path, method, headers)
	if !isOk {
		return statusCode, user
	}

	switch path[1:5] {
	case "user":
		return ProcesoUsers(body, path, method, user, id, request)

	case "prod":
		return ProcesoProducts(body, path, method, user, idn, request)

	case "stoc":
		return ProcesoStocks(body, path, method, user, idn, request)

	case "addr":
		return ProcesoAddress(body, path, method, user, idn, request)

	case "cate":
		return ProcesoCategory(body, path, method, user, idn, request)

	case "orde":
		return ProcesoOrders(body, path, method, user, idn, request)

	}

	return 400, "Method Invalid SALAME"

}

func validoAuthorization(path string, method string, headers map[string]string) (bool, int, string) {
	if (path == "product" && method == "GET") ||
		(path == "category" && method == "GET") {
		return true, 200, ""
	}

	token := headers["authorization"]
	if len(token) == 0 {
		return false, 401, "Token requerido"
	}

	todoOK, err, msg := auth.ValidoToken(token)

	if !todoOK {
		if err != nil {
			fmt.Println("Error en el token " + err.Error())
			return false, 401, err.Error()
		} else {
			fmt.Println("Error en el token " + msg)
			return false, 401, msg
		}
	}
	fmt.Println("Token OK")
	return true, 200, msg
}

func ProcesoUsers(body string, path string, method string, user string, id string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Ingreando a la funcion ProcesoUsers")
	fmt.Println("Imprimo Path")
	fmt.Println(path)

	if path == "/user/me" {
		switch method {
		case "PUT":
			return routers.UpdateUser(body, user)
		case "GET":
			return routers.SelectUser(body, user)
		}
	}
	if path == "/users" {
		switch method {
		case "GET":
			return routers.SelectUsers(body, user, request)
		}

	}

	return 400, "method invalid ProcesoUsers"
}

func ProcesoProducts(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Ingreando a la funcion ProcesoProducts")

	switch method {
	case "POST":
		return routers.InsertProduct(body, user)
	case "PUT":
		return routers.UpdateProduct(body, user, id)
	case "DELETE":
		return routers.DeleteProduct(user, id)
	case "GET":
		return routers.SelectProduct(body, request)
	}

	return 400, "method invalid"
}

func ProcesoCategory(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	switch method {
	case "POST":
		return routers.InsertCategory(body, user)
	case "PUT":
		return routers.UpdateCategory(body, user, id)
	case "DELETE":
		return routers.DeleteCategory(body, user, id)
	case "GET":
		return routers.SelectCategories(body, request)
	}
	return 400, "Method Invalid ProcesoCategory"
}

func ProcesoStocks(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Ingreando a la funcion ProcesoStocks")

	return routers.UpdateStock(body, user, id)
}

func ProcesoAddress(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Ingreando a la funcion ProcesoAddress")

	switch method {
	case "POST":
		return routers.InsertAddress(body, user)
	case "PUT":
		return routers.UpdateAddress(body, user, id)
	case "DELETE":
		return routers.DeleteAddres(user, id)
	case "GET":
		return routers.SelectAddress(user)
	}

	return 400, "method invalid"
}

func ProcesoOrders(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Ingreando a la funcion ProcesoOrder")

	switch method {
	case "POST":
		return routers.InsertOrder(body, user)
	case "GET":
		return routers.SelectOrders(user, request)

	}

	return 400, "method invalid"
}
