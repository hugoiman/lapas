package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// MySigningKey is signature
var MySigningKey = []byte("jwt super secret key")

// MyClaims is Credential
type MyClaims struct {
	IDUser  int    `json:"idUser"`
	Nama    string `json:"nama"`
	Job     string `json:"job"`
	Pangkat string `json:"pangkat"`
	Divisi  string `json:"divisi"`
	jwt.StandardClaims
}

// AuthToken is middleware
func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			http.Error(w, "Gagal! Dibutuhkan otentikasi. Silahkan melakukan login.", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		claims := &MyClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return MySigningKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Token expired/key tidak cocok(invalid)
			return
		}
		if !token.Valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		context.Set(r, "user", claims)

		next.ServeHTTP(w, r)
	})
}
