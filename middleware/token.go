package middleware

import (
	"context"
	"lapas/pkg/controllers"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthToken is middleware
func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			http.Error(w, "Gagal! Dibutuhkan otentikasi. Silahkan melakukan login.", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		claims := &controllers.MyClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return controllers.MySigningKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Token expired/key tidak cocok(invalid)
			return
		}
		if !token.Valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(context.Background(), userInfo, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type key string

var userInfo = key("user")
