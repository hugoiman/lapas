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
		if user.Divisi != "SDM" {
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
		// allowed -> Direksi/Direktur/Sekretaris Perusahaan/Admin
		if user.Job != "Direktur" && user.Job != "Direksi" && user.Divisi != "Sekretaris Perusahaan" && user.Divisi != "Logistik & Administrasi" {
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
		if user.Divisi != "Sekretaris Perusahaan" && user.Divisi == "Logistik & Administrasi" {
			http.Error(w, "Gagal! Anda tidak diizinkan.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// CDispo is middleware
func CDispo(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*MyClaims)
		// allowed -> Direksi/Direktur/Sekretaris Perusahaan/Kadiv/Manager/Assistant Manager
		if user.Divisi != "Sekretaris Perusahaan" && !(user.Job == "Direksi" || user.Job == "Direktur" || user.Pangkat == "Assistant Vice President" || user.Pangkat == "Manager" || user.Pangkat == "Assistant Manager") {
			http.Error(w, "Gagal! Anda tidak diizinkan.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
