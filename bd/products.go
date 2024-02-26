package bd

import (
	"database/sql"
	"errors"

	//"errors"
	"fmt"
	"strconv"

	"strings"

	"github.com/gambit/models"
	"github.com/gambit/tools"
	_ "github.com/go-sql-driver/mysql"
)

func InsertProduct(p models.Product) (int64, error) {
	fmt.Println("Inicializando funcion  db.InsertProduct")

	err := DbConnect()
	if err != nil {
		return 0, err
	}

	defer Db.Close()

	sentencia := "INSERT INTO products (Prod_Title"

	if len(p.ProdDescription) > 0 {
		sentencia += ", Prod_Description"
	}

	if p.ProdPrice > 0 {
		sentencia += ", Prod_Price"
	}

	if p.ProdCategId > 0 {
		sentencia += ", Prod_CategoryId"
	}

	if p.ProdStock > 0 {
		sentencia += ", Prod_Stock"
	}

	if len(p.ProdPath) > 0 {
		sentencia += ", Prod_Path"
	}

	sentencia += ") Values ('" + tools.EscapeString(p.ProdTitle) + "'"

	if len(p.ProdDescription) > 0 {
		sentencia += ", '" + tools.EscapeString(p.ProdDescription) + "'"
	}

	if p.ProdPrice > 0 {
		sentencia += ", " + strconv.FormatFloat(p.ProdPrice, 'e', -1, 64)
	}

	if p.ProdCategId > 0 {
		sentencia += ", " + strconv.Itoa(p.ProdCategId)
	}

	if p.ProdStock > 0 {
		sentencia += ", " + strconv.Itoa(p.ProdStock)
	}

	if len(p.ProdPath) > 0 {
		sentencia += ", '" + tools.EscapeString(p.ProdPath) + "'"
	}

	sentencia += ")"

	var result sql.Result

	fmt.Println("Vamos a ejecutar la sentancia: ")
	fmt.Println(sentencia)
	result, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	LastInsertId, err2 := result.LastInsertId()
	if err2 != nil {
		fmt.Println(err.Error())
		return 0, err2
	}

	fmt.Println("Insert Product > Ejecucion Exitosa")
	return LastInsertId, nil
}

func UpdateProduct(p models.Product) error {
	fmt.Println("Inicializando funcion  db.UpdateProduct")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "Update products SET "

	sentencia = tools.ArmoSentencia(sentencia, "Prod_Title", "S", 0, 0, p.ProdTitle)
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Description", "S", 0, 0, p.ProdDescription)
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Price", "F", 0, p.ProdPrice, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_CategoryId", "N", p.ProdCategId, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Stock", "N", p.ProdStock, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Path", "S", 0, 0, p.ProdPath)

	sentencia += " Where Prod_Id = " + strconv.Itoa(p.ProdId)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Product > ejecutado Exitosamente")
	return nil
}

func DeleteProduct(id int) error {
	fmt.Println("Inicializando funcion  db.DeleteProduct")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "Delete from products Where Prod_Id =" + strconv.Itoa(id)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Product > ejecutado Exitosamente")
	return nil
}

func SelectProduct(p models.Product, choice string, page int, pageSize int, orderType string, orderField string) (models.ProductResp, error) {
	fmt.Println("Inicializando funcion  db.SelectProduct")

	var Resp models.ProductResp
	var Prod []models.Product

	err := DbConnect()
	if err != nil {
		return Resp, err
	}

	defer Db.Close()

	//Primera sentencia para contar cuantos productos existen
	var sentencia string
	var sentenciaCount string
	var where, limit string

	sentencia = "Select Prod_Id, Prod_Title, Prod_Description, Prod_CreatedAt, Prod_Updated, Prod_Price, Prod_Path, Prod_CategoryId, Prod_Stock FROM  products "
	sentenciaCount = "Select count(*) as registros from products "

	switch choice {
	case "P":
		where = " Where Prod_Id = " + strconv.Itoa(p.ProdId)
	case "S":
		where = " Where UCASE(CONCAT(Prod_Title, Prod_Description)) LIKE '%" + strings.ToUpper(p.ProdSearch) + "%' "
	case "C":
		where = " Where Prod_CategoryId = " + strconv.Itoa(p.ProdCategId)
	case "U":
		where = " Where UCASE(Prod_Path) LIKE '%" + strings.ToUpper(p.ProdPath) + "%' "
	case "K":
		join := " JOIN category on Prod_CategoryId = Categ_Id AND Categ_Path like '%" + strings.ToUpper(p.ProdCategPath) + "%' "
		sentencia += join
		sentenciaCount += join

	}

	sentenciaCount += where

	fmt.Println("Valor de sentenciaCount")
	fmt.Println(sentenciaCount)

	var rows *sql.Rows
	rows, err = Db.Query(sentenciaCount)

	if err != nil {
		fmt.Println(err.Error())
		return Resp, err
	}
	//Nos aseguramos que cierra la sentencia si se dio error
	defer rows.Close()

	//Colocamos el valor de la sentencia en la variable creada.
	rows.Next()
	var regi sql.NullInt32
	//err = rows.Scan(&regi)

	registros := int(regi.Int32)

	if page > 0 {
		if registros > pageSize {
			limit = " LIMIT " + strconv.Itoa(pageSize)
			if page > 1 {
				offset := pageSize * (page - 1)
				limit += " OFFSET " + strconv.Itoa(offset)
			}
		} else {
			limit = ""
		}
	}

	var orderBy string
	if len(orderField) > 0 {
		switch orderField {
		case "I":
			orderBy = " ORDER BY Prod_Id "
		case "T":
			orderBy = " ORDER BY Prod_Title "
		case "D":
			orderBy = " ORDER BY Prod_Description "
		case "F":
			orderBy = " ORDER BY Prod_CreatedAt "
		case "P":
			orderBy = " ORDER BY Prod_Price "
		case "S":
			orderBy = " ORDER BY Prod_Stock "
		case "C":
			orderBy = " ORDER BY Prod_Id "

		}

		if orderType == "D" {
			orderBy += " DESC "
		}
	}

	sentencia += where + orderBy + limit
	fmt.Println("Valor de sentencia")
	fmt.Println(sentencia)

	rows, _ = Db.Query(sentencia)
	for rows.Next() {
		var pp models.Product
		var ProdId sql.NullInt32
		var ProdTitle sql.NullString
		var ProdDescription sql.NullString
		var ProdCreatedAt sql.NullString
		var ProdUpdated sql.NullString
		var ProdPrice sql.NullFloat64
		var ProdPath sql.NullString
		var ProdCategoryId sql.NullInt32
		var ProdStock sql.NullInt32

		err := rows.Scan(&ProdId, &ProdTitle, &ProdDescription, &ProdCreatedAt, &ProdUpdated, &ProdPrice, &ProdPath, &ProdCategoryId, &ProdStock)

		if err != nil {
			return Resp, err
		}

		pp.ProdId = int(ProdId.Int32)
		pp.ProdTitle = ProdTitle.String
		pp.ProdDescription = ProdDescription.String
		pp.ProdCreatedAt = ProdCreatedAt.String
		pp.ProdUpdated = ProdUpdated.String
		pp.ProdPrice = ProdPrice.Float64
		pp.ProdPath = ProdPath.String
		pp.ProdCategId = int(ProdCategoryId.Int32)
		pp.ProdStock = int(ProdStock.Int32)
		Prod = append(Prod, pp)
	}

	Resp.TotalItems = registros
	Resp.Data = Prod

	fmt.Println("Select Product > Ejecucion exitosa")
	return Resp, nil
}

func UpdateStock(p models.Product) error {
	fmt.Println("Inicializando funcion  db.UpdateStock")

	if p.ProdStock == 0 {
		return errors.New("[ERROR] debe enviar el stock a modificar")
	}

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "Update products  SET Prod_Stock = Prod_Stock  +" + strconv.Itoa(p.ProdStock) + " Where Prod_Id = " + strconv.Itoa(p.ProdId)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Stock > ejecutado Exitosamente")
	return nil
}
