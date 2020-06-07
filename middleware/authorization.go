package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

// IsIT is middleware
func IsIT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*MyClaims)
		if !strings.EqualFold(user.Divisi, "IT") {
			http.Error(w, "Gagal! Anda bukan Divisi IT.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// IsSDM is middleware
func IsSDM(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*MyClaims)
		if !strings.EqualFold(user.Divisi, "SDM") {
			http.Error(w, "Gagal! Anda bukan Divisi SDM.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RSurat is middleware
func RSurat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*MyClaims)
		// fmt.Printf("%+v", user)
		if !strings.EqualFold(user.Job, "Direksi") {
			http.Error(w, "Gagal! Anda tidak diizinkan.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// CUDSurat is middleware
func CUDSurat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*MyClaims)
		if !strings.EqualFold(user.Divisi, "Sekretaris Perusahaan") {
			http.Error(w, "Gagal! Anda tidak diizinkan.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
