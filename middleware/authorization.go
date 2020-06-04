package middleware

import (
	"lapas/pkg/controllers"
	"net/http"
	"strings"
)

// IsIT is middleware
func IsIT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userInfo).(*controllers.MyClaims)
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
		user := r.Context().Value(userInfo).(*controllers.MyClaims)
		if !strings.EqualFold(user.Divisi, "SDM") {
			http.Error(w, "Gagal! Anda bukan Divisi SDM.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
