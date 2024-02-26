package bd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	//"strconv"
	"github.com/gambit/models"
	"github.com/gambit/tools"
	_ "github.com/go-sql-driver/mysql"
)

func InsertCategory(c models.Category) (int64, error) {
	fmt.Println("Inicializando funcion  db.InsertCategory")

	err := DbConnect()
	if err != nil {
		return 0, err
	}

	defer Db.Close()

	sentencia := "INSERT INTO category (Categ_Name, Categ_Path ) Values('" + c.CategName + "','" + c.CategPath + "')"

	fmt.Println("ejecutnado sentencia")
	fmt.Println(sentencia)
	var result sql.Result

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

	fmt.Println("Insert Category > Ejecucion Exitosa")
	return LastInsertId, nil

}

func UpdateCategory(c models.Category) error {
	fmt.Println("Comienza registro de UpdateCategory")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "UPDATE category SET "

	if len(c.CategName) > 0 {
		sentencia += " Categ_Name = '" + tools.EscapeString(c.CategName) + "'"
	}

	if len(c.CategPath) > 0 {
		if !strings.HasSuffix(sentencia, "SET ") {
			sentencia += ", "
		}
		sentencia += "Categ_Path = '" + tools.EscapeString(c.CategPath) + "'"

	}

	sentencia += " WHERE Categ_Id = " + strconv.Itoa(c.CategID)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Category > Ejecucion Exitosa")
	return nil
}

func DeleteCategory(id int) error {
	fmt.Println("Comienza registro de DeleteCategory")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "DELETE FROM category WHERE Categ_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Category > Ejecucion Exitosa")
	return nil

}

func SelectCategories(CategId int, Slug string) ([]models.Category, error) {
	fmt.Println("Inicializando funcion  db.SelectCategories")

	var Categ []models.Category

	err := DbConnect()

	if err != nil {
		return Categ, err
	}

	defer Db.Close()

	sentencia := "Select Categ_Id, Categ_Name, Categ_Path  from category  "

	if CategId > 0 {
		sentencia += " Where Categ_Id = " + strconv.Itoa(CategId)
	} else {
		if len(Slug) > 0 {
			sentencia += " Where Categ_Path like '%" + Slug + "%'"
		}
	}

	fmt.Println(sentencia)
	var rows *sql.Rows

	rows, _ = Db.Query(sentencia)

	//Nos movemos entre las filas para colocar los registros
	for rows.Next() {
		//Definimos variables, evitamos dejarlos en nulos
		var c models.Category
		var categId sql.NullInt32
		var categName sql.NullString
		var categPath sql.NullString

		//Colocamos los varores  de los registros en las variebles
		err := rows.Scan(&categId, &categName, &categPath)
		if err != nil {
			return Categ, err
		}

		//Colocamos los registros en la variable de modelo.
		c.CategID = int(categId.Int32)
		c.CategName = categName.String
		c.CategPath = categPath.String

		//Sumamos los registros al slice
		Categ = append(Categ, c)

	}

	fmt.Println("Select Category > Ejecucion exitosa")

	return Categ, nil

}
