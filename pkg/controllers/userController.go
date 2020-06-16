package controllers

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	models "lapas/pkg/models"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/gomail.v2"
)

// GetUser is function
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idUser := vars["idUser"]
	data := models.GetUser(idUser)
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetUsers is function
func GetUsers(w http.ResponseWriter, r *http.Request) {
	data := models.GetUsers()
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateUser is New User
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Printf("%+v\n", user)

	if err := validator.New().Struct(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// UpdateUser is Edit User
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	vars := mux.Vars(r)
	idUser := vars["idUser"]

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.UpdateUser(idUser, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil diperbarui!"}`))
}

// ChangePassword is Edit Password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idUser := vars["idUser"]

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	if err := validator.New().Var(fmt.Sprintf("%v", data["password_baru"]), "required,min=6,max=18"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Var(fmt.Sprintf("%v", data["password_lama"]), "required,min=6,max=18"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var oldPass = sha1.New()
	oldPass.Write([]byte(fmt.Sprintf("%v", data["password_lama"])))
	var encryptedOldPass = fmt.Sprintf("%x", oldPass.Sum(nil))

	isValid := models.CheckOldPassword(idUser, encryptedOldPass)
	if !isValid {
		http.Error(w, "Password lama tidak sesuai", http.StatusBadRequest)
		return
	}

	var newPass = sha1.New()
	newPass.Write([]byte(fmt.Sprintf("%v", data["password_baru"])))
	var encryptedPass = fmt.Sprintf("%x", newPass.Sum(nil))

	models.UpdatePassword(idUser, encryptedPass)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Password berhasil diperbarui!"}`))
}

// ResetPassword is forgot password
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	nipg := fmt.Sprintf("%v", data["nipg"])
	email := fmt.Sprintf("%v", data["email"])

	if err := validator.New().Var(email, "required,email"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Var(nipg, "required"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := models.CheckUser(nipg, email)
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusBadRequest)
		return
	}

	// Generate Random String
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randomStr := make([]rune, 10)
	for i := range randomStr {
		randomStr[i] = letters[rand.Intn(len(letters))]
	}
	newPass := string(randomStr)

	var pass = sha1.New()
	pass.Write([]byte(newPass))
	var encryptedPass = fmt.Sprintf("%x", pass.Sum(nil))
	fmt.Println(encryptedPass)

	// models.UpdatePassword(idUser, encryptedPass)
	// sendNotifikasi(newPass, email, message)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Password baru telah dikirim ke email pengguna."}`))
}

// SendEmail is func
func SendEmail(subject, to, message string) {
	var configSMTPHost = "smtp.gmail.com"
	var configSMTPPort = 587
	var configEmail = "nanonymoux@gmail.com"
	var configPassword = "kudaponi10"

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", configEmail)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetHeader("text/html", message)

	dialer := gomail.NewDialer(configSMTPHost, configSMTPPort, configEmail, configPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("errpor :", err.Error())
	}

}
