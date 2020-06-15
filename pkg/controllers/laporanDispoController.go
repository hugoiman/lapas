package controllers

import (
	"encoding/json"
	models "lapas/pkg/models"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// CreateLaporanDispo is func
func CreateLaporanDispo(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idLaporan := vars["idLaporan"]

	var laporan models.LaporanDispo

	if err := json.NewDecoder(r.Body).Decode(&laporan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(laporan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getLaporan, err := models.GetLaporanDispo(idLaporan)
	if err != nil {
		http.Error(w, "Gagal! Laporan disposisi tidak ditemukan.", http.StatusBadRequest)
		return
	} else if getLaporan.IDPenerima != user.IDUser {
		http.Error(w, "Gagal! Anda bukan penerima disposisi.", http.StatusBadRequest)
		return
	}

	laporan.IDDisposisi = getLaporan.IDDisposisi
	laporan.Status = "Report"
	laporan.UpdatedAt = time.Now().Format("2006-01-02")

	models.CreateLaporanDisposisi(idLaporan, laporan)
}
