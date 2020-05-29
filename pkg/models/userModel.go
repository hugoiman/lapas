package models

import (
	"fmt"
	"lapas/db"
	"time"
)

// User Class
type User struct {
	IDUser     int       `json:"idUser"`
	NIPG       string    `json:"nipg" validate:"required"`
	Nama       string    `json:"nama" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Pangkat    string    `json:"pangkat" validate:"required"`
	Divisi     string    `json:"divisi" validate:"required"`
	Direktorat string    `json:"direktorat" validate:"required"`
	Actived    bool      `json:"actived" validate:"required"`
	TglLahir   time.Time `json:"tglLahir" validate:"required"`
}

// Users is User List
type Users struct {
	Users []User `json:"user"`
}

// GetUser is function
func GetUser(idUser string) User {
	con := db.Connect()
	query := "SELECT idUser, nipg, nama, email, pangkat, divisi, direktorat, actived from user where idUser = ?"

	user := User{}
	err := con.QueryRow(query, idUser).Scan(
		&user.IDUser, &user.NIPG, &user.Nama, &user.Email,
		&user.Pangkat, &user.Direktorat, &user.Divisi, &user.Actived)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer con.Close()

	return user
}

// GetUsers is function
func GetUsers() Users {
	con := db.Connect()
	query := "SELECT idUser, nipg, nama, email, pangkat, divisi, direktorat, actived, tglLahir from user LIMIT 10"
	rows, err := con.Query(query)

	if err != nil {
		fmt.Println(err.Error())
	}

	user := User{}
	users := Users{}

	for rows.Next() {
		err2 := rows.Scan(
			&user.IDUser, &user.NIPG, &user.Nama, &user.Email,
			&user.Pangkat, &user.Divisi, &user.Direktorat, &user.Actived, &user.TglLahir)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
		users.Users = append(users.Users, user)
	}

	defer con.Close()
	return users

}

// CreateUser is New User
func CreateUser(user User) error {
	con := db.Connect()
	_, err := con.Exec("INSERT INTO user (nipg, nama, email, pangkat, divisi, direktorat, actived, tglLahir) VALUES (?,?,?,?,?,?,?,?)", user.NIPG, user.Nama, user.Email, user.Pangkat, user.Divisi, user.Direktorat, user.Actived, user.TglLahir)

	defer con.Close()

	return err
}

// UpdateUser is Edit User
func UpdateUser(idUser string, user User) error {
	con := db.Connect()
	query := "UPDATE user SET nipg = ?, nama = ?, email = ?, pangkat = ?, divisi = ?, direktorat = ?, actived = ?, tglLahir = ? WHERE idUser = ?"
	_, err := con.Exec(query, user.NIPG, user.Nama, user.Email, user.Pangkat, user.Divisi, user.Direktorat, user.Actived, user.TglLahir, idUser)

	defer con.Close()

	return err
}

// CheckOldPassword is Auth User
func CheckOldPassword(idUser, password string) bool {
	var isAny bool
	con := db.Connect()
	query := "SELECT EXISTS (SELECT 1 FROM `user` WHERE idUser = ? AND password = ?)"
	con.QueryRow(query, idUser, password).Scan(&isAny)

	defer con.Close()

	return isAny
}

// UpdatePassword is Edit Password
func UpdatePassword(idUser, password string) {
	con := db.Connect()
	query := "UPDATE user SET password = ? WHERE idUser = ?"
	con.Exec(query, password, idUser)

	defer con.Close()
}

// CheckUser is function
func CheckUser(nipg, email string) (string, error) {
	var idUser string
	con := db.Connect()
	query := "SELECT idUser FROM user WHERE nipg = ? AND email = ?"
	err := con.QueryRow(query, nipg, email).Scan(&idUser)

	defer con.Close()
	return idUser, err
}
