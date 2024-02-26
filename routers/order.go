package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gambit/bd"
	"github.com/gambit/models"
)

func InsertOrder(body string, User string) (int, string) {
	fmt.Println("Inicializando funcion  router.InsertOrder")

	var o models.Orders

	err := json.Unmarshal([]byte(body), &o)
	if err != nil {
		return 400, "Error en los datos recibidos" + err.Error()
	}

	o.Order_UserUUID = User

	OK, message := ValidOrder(o)
	if !OK {
		return 400, message
	}

	result, err2 := bd.InsertOrder(o)
	if err2 != nil {
		return 400, "Ocurrio un error al intentar realizar el registro de la orden " + err2.Error()
	}

	return 200, " { Order ID :" + strconv.Itoa(int(result)) + "}"

}

func ValidOrder(o models.Orders) (bool, string) {
	fmt.Println("Inicializando funcion  router.ValidOrder")

	if o.Order_Total == 0 {
		return false, "Debe indicar un total de orden"
	}

	count := 0
	for _, od := range o.OrdersDetails {
		if od.OD_ProdId == 0 {
			return false, "Debe indicar el id en el detalle de la orden"
		}
		if od.OD_Quantity == 0 {
			return false, "Debe indicar la cantidad del producto en el detalle de la orden"
		}
		count++
	}

	if count == 0 {
		return false, "Debe indicar items en la orden"
	}
	return true, ""
}

func SelectOrders(user string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Inicializando funcion  router.SelectOrders")

	var fechaDesde, fechaHasta string
	var orderId int
	var page int

	if len(request.QueryStringParameters["fechaDesde"]) > 0 {
		fechaDesde = request.QueryStringParameters["fechaDesde"]
	}

	if len(request.QueryStringParameters["fechaHasta"]) > 0 {
		fechaHasta = request.QueryStringParameters["fechaHasta"]
	}

	if len(request.QueryStringParameters["page"]) > 0 {
		page, _ = strconv.Atoi(request.QueryStringParameters["page"])
	}
	if len(request.QueryStringParameters["orderId"]) > 0 {
		orderId, _ = strconv.Atoi(request.QueryStringParameters["orderId"])
	}

	result, err := bd.SelectOrders(user, fechaDesde, fechaHasta, page, orderId)
	if err != nil {
		return 400, "Ocurri un error al intentar capturar los registros de ordenes del " + fechaDesde + " al " + fechaHasta + " > " + err.Error()
	}

	orders, err2 := json.Marshal(result)
	if err2 != nil {
		return 400, "Ocurrio un error al intentar convertir en JSON  la busqueda de ordenes/s > " + err2.Error()
	}

	return 200, string(orders)

}
