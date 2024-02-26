package bd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gambit/models"
	_ "github.com/go-sql-driver/mysql"
)

func InsertAddress(addr models.Address, User string) error {
	fmt.Println("Inicializando funcion  bd.InsertAddress")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "INSERT INTO addresses (Add_UserId, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name )"
	sentencia += " VALUES ('" + User + "','" + addr.AddAddress + "','" + addr.AddCity + "','" + addr.AddState + "','"
	sentencia += addr.AddPostalCode + "','" + addr.AddPhone + "','" + addr.AddName + "','" + addr.AddName + "')"

	fmt.Println("Vamos a ejecutar la sentancia: ")
	fmt.Println(sentencia)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Insert address > Ejecucion exitosa ")
	return nil
}

func AddressExists(User string, id int) (error, bool) {
	fmt.Println("Inicializando funcion  bd.AddessExists")

	err := DbConnect()
	if err != nil {
		return err, false
	}

	defer Db.Close()

	sentencia := "Select 1 from addresses Where add_Id = " + strconv.Itoa(id) + " AND Add_UserId = '" + User + "'"
	fmt.Println(sentencia)

	rows, err := Db.Query(sentencia)
	if err != nil {
		return err, false
	}

	var valor string
	rows.Next()
	rows.Scan(&valor)

	fmt.Println("AddressExist > Ejecucion exitosa - valor devuelto  " + valor)

	if valor == "1" {
		return nil, true
	}
	return nil, false

}

func UpdateAddress(addr models.Address) error {
	fmt.Println("Inicializando funcion  bd.UpdateAddress")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "Update addresses SET "

	if addr.AddAddress != "" {
		sentencia += "Add_Address = '" + addr.AddAddress + "', "
	}

	if addr.AddCity != "" {
		sentencia += "Add_City = '" + addr.AddCity + "', "
	}

	if addr.AddName != "" {
		sentencia += "Add_Name = '" + addr.AddName + "', "
	}

	if addr.AddPhone != "" {
		sentencia += "Add_Phone = '" + addr.AddPhone + "', "
	}

	if addr.AddPostalCode != "" {
		sentencia += "Add_PostalCode = '" + addr.AddPostalCode + "', "
	}

	if addr.AddState != "" {
		sentencia += "Add_State = '" + addr.AddState + "', "
	}

	if addr.AddTitle != "" {
		sentencia += "Add_Title = '" + addr.AddTitle + "', "
	}

	sentencia, _ = strings.CutSuffix(sentencia, ", ")
	sentencia += " Where Add_Id = " + strconv.Itoa(addr.AddId)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Address > ejecutado Exitosamente")
	return nil

}

func DeleteAddress(id int) error {
	fmt.Println("Inicializando funcion  db.DeleteAddress")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "Delete from addresses Where Add_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Address > ejecutado Exitosamente")
	return nil
}

func SelectAddress(User string) ([]models.Address, error) {
	fmt.Println("Inicializando funcion  db.SelectAddress")

	address := []models.Address{}

	err := DbConnect()
	if err != nil {
		return address, err
	}

	defer Db.Close()

	sentencia := "Select Add_Id, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name From addresses Where Add_UserId = '" + User + "'"

	fmt.Println("Valor de sentencia")
	fmt.Println(sentencia)

	var rows *sql.Rows
	rows, err = Db.Query(sentencia)

	if err != nil {
		fmt.Println(err.Error())
		return address, err
	}

	//Nos aseguramos que cierra la sentencia
	defer rows.Close()

	for rows.Next() {
		var a models.Address
		var addId sql.NullInt16
		var addAddress sql.NullString
		var addCity sql.NullString
		var addState sql.NullString
		var addPostalCode sql.NullString
		var addPhone sql.NullString
		var addTitle sql.NullString
		var addName sql.NullString

		err := rows.Scan(&addId, &addAddress, &addCity, &addState, &addPostalCode, &addPhone, &addTitle, &addName)
		if err != nil {
			return address, err
		}

		a.AddId = int(addId.Int16)
		a.AddAddress = addAddress.String
		a.AddCity = addAddress.String
		a.AddState = addState.String
		a.AddPostalCode = addPostalCode.String
		a.AddPhone = addPhone.String
		a.AddTitle = addTitle.String
		a.AddName = addName.String
		address = append(address, a)

	}

	fmt.Println("Select Address > Ejecucion Exitosa ")
	return address, nil

}
