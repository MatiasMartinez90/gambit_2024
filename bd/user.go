package bd

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gambit/models"
	"github.com/gambit/tools"
	_ "github.com/go-sql-driver/mysql"
)

func UpdateUser(UField models.User, User string) error {
	fmt.Println("Inicializando funcion  db.UpdateUser")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	sentencia := "UPDATE users SET "

	coma := ""
	if len(UField.UserFirstName) > 0 {
		coma = ","
		sentencia += "User_FirstName = '" + UField.UserFirstName + "'"
	}

	if len(UField.UserLastName) > 0 {
		sentencia += coma + "User_LastName = '" + UField.UserLastName + "'"
	}

	sentencia += ", User_DateUpg = '" + tools.FechaMySQL() + "' Where User_UUID='" + User + "'"

	println("ejecutando sentencia")
	println(sentencia)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update de User > Ejecucion exitosa")
	return nil

}

func SelectUser(UserId string) (models.User, error) {
	fmt.Println("Inicializando funcion  db.SelectUser")

	User := models.User{}

	err := DbConnect()
	if err != nil {
		return User, err
	}

	defer Db.Close()

	sentencia := "Select * from users Where User_UUID = '" + UserId + "'"

	var rows *sql.Rows
	rows, err = Db.Query(sentencia)
	//Nos aseguramos que cierra la sentencia

	if err != nil {
		fmt.Println(err.Error())
		return User, err
	}

	defer rows.Close()

	rows.Next()

	var firstName sql.NullString
	var lastName sql.NullString
	var dateUpg sql.NullString

	err = rows.Scan(&User.UserUUID, &User.UserEmail, &firstName, &lastName, &User.UserStatus, &User.UserDateAdd, &dateUpg)

	fmt.Println(err)

	User.UserFirstName = firstName.String
	User.UserLastName = lastName.String
	User.UserDateUpd = dateUpg.String

	fmt.Println("Select User > Ejecucion exitosa")
	return User, nil

}

func SelectUsers(Page int) (models.ListUsers, error) {
	fmt.Println("Inicializando funcion  db.SelectUsers")

	var listUser models.ListUsers
	User := []models.User{}

	err := DbConnect()
	if err != nil {
		return listUser, err
	}

	defer Db.Close()

	var offset int = (Page * 10) - 10
	var sentencia string
	var sentenciaCount string = "Select count(*) as registros from users"

	sentencia = "Select User_UUID, User_Email, User_FirstName, User_LastName, User_Status, User_DateAdd, User_DateUpg from users limit 10 "
	if offset > 0 {
		sentencia = " OFFSET " + strconv.Itoa(offset)
	}

	var rowsCount *sql.Rows

	fmt.Print("Ejecutando sentenciaCount")
	fmt.Println(sentenciaCount)
	rowsCount, err = Db.Query(sentenciaCount)
	if err != nil {
		return listUser, err
	}

	defer rowsCount.Close()

	rowsCount.Next()

	var registros int
	rowsCount.Scan(&registros)
	listUser.TotalItems = registros

	fmt.Print("Ejecutando sentencia simple")
	fmt.Println(sentencia)

	var rows *sql.Rows

	rows, err = Db.Query(sentencia)
	if err != nil {
		return listUser, err
	}

	for rows.Next() {
		var u models.User

		var firstName sql.NullString
		var lastName sql.NullString
		var dateUpg sql.NullString

		err = rows.Scan(&u.UserUUID, &u.UserEmail, &firstName, &lastName, &u.UserStatus, &u.UserDateAdd, &dateUpg)

		fmt.Println(err)

		u.UserFirstName = firstName.String
		u.UserLastName = lastName.String
		u.UserDateUpd = dateUpg.String
		User = append(User, u)
	}

	fmt.Println("Select Users > Ejecucion exitosa")
	listUser.Data = User
	return listUser, nil

}
