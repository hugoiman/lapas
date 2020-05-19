package models

import (
	"fmt"
	"lapas/db"
	"time"
)

// User Class
type User struct {
	IDUser     int       `json:"idUser"`
	NIPG       string    `json:"nipg"`
	Nama       string    `json:"nama"`
	Email      string    `json:"email"`
	Pangkat    string    `json:"pangkat"`
	Direktorat string    `json:"direktorat"`
	Divisi     string    `json:"divisi"`
	Actived    bool      `json:"activated"`
	TglLahir   time.Time `json:"tglLahir"`
}

// GetUser is function
func GetUser(idUser string) User {
	con := db.Connect()
	query := "SELECT idUser, nipg, nama, email, pangkat, direktorat, divisi, actived, tglLahir from user where idUser = ?"

	user := User{}
	err := con.QueryRow(query, idUser).Scan(
		&user.IDUser, &user.NIPG, &user.Nama, &user.Email,
		&user.Pangkat, &user.Direktorat, &user.Divisi, &user.Actived, &user.TglLahir)

	if err != nil {
		fmt.Print(err.Error())
		fmt.Println("Error")
	}

	defer con.Close()

	return user
}

// GetUsers is function
func GetUsers() []User {
	con := db.Connect()
	query := "SELECT idUser, nipg, nama, email, pangkat, direktorat, divisi, actived, tglLahir from user LIMIT 10"
	rows, err := con.Query(query)

	if err != nil {
		fmt.Println(err.Error())
	}

	user := User{}
	userList := []User{}

	for rows.Next() {
		err2 := rows.Scan(
			&user.IDUser, &user.NIPG, &user.Nama, &user.Email,
			&user.Pangkat, &user.Direktorat, &user.Divisi, &user.Actived, &user.TglLahir)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
		userList = append(userList, user)
	}

	defer con.Close()
	return userList

}
