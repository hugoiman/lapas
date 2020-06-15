package controllers

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	mw "lapas/middleware"
	models "lapas/pkg/models"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/gomail.v2"
)

// Login is func
func Login(w http.ResponseWriter, r *http.Request) {
	var login models.Auth
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sha = sha1.New()
	sha.Write([]byte(login.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	login.Password = encryptedString

	idUser, err := models.Login(login)
	if err != nil {
		http.Error(w, "Gagal! Email/nipg atau password salah.", http.StatusBadRequest)
		return
	}

	user := models.GetUser(idUser)
	token := CreateToken(user)

	type M map[string]interface{}
	message, _ := json.Marshal(M{"token": token})

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// MySigningKey is signature
var MySigningKey = mw.MySigningKey

// MyClaims is credential
type MyClaims = mw.MyClaims

// CreateToken is Generate token
func CreateToken(user models.User) string {
	claims := MyClaims{
		IDUser:  user.IDUser,
		Nama:    user.Nama,
		Job:     user.Job,
		Pangkat: user.Pangkat,
		Divisi:  user.Divisi,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(MySigningKey)

	return tokenString
}

// Logout is func
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "tokenCookie",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "localhost:8000/", http.StatusSeeOther)
	return
}

const configSMTPHost = "smtp.gmail.com"
const configSMTPPort = 587
const configEmail = "nanonymoux@gmail.com"
const configPassword = "kudaponi10"

// SendEmail is func
func SendEmail(subject, to, message string) {
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
