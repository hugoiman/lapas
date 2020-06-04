package models

import "lapas/db"

// Auth is class
type Auth struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login is func
func Login(login Auth) (string, error) {
	var idUser string

	con := db.Connect()
	query := "SELECT idUser FROM user WHERE (nipg = ? OR email = ?) AND password = ? AND actived = 1"
	err := con.QueryRow(query, login.ID, login.ID, login.Password).Scan(&idUser)

	defer con.Close()
	return idUser, err

}
