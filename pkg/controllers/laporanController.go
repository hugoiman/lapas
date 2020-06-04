package controllers

import (
	"encoding/json"
	"fmt"
	models "lapas/pkg/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetLaporan is func
func GetLaporan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idLaporan := vars["idLaporan"]

	data, err := models.GetLaporan(idLaporan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Laporan Not Found
		return
	}

	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetMyLaporan is get list sub
func GetMyLaporan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idUser := vars["idUser"]

	data := models.GetMyLaporan(idUser)
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetLaporans is get list sub
func GetLaporans(w http.ResponseWriter, r *http.Request) {
	data := models.GetLaporans()
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateLaporan is func
func CreateLaporan(w http.ResponseWriter, r *http.Request) {
	var laporan models.Laporan

	if err := json.NewDecoder(r.Body).Decode(&laporan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(laporan); err != nil {
		fmt.Println()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	laporan.TglKirim = time.Now().Format("2006-01-02")
	laporan.Status = "Terkirim"

	err := models.CreateLaporan(laporan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// CreateTanggapan is func
func CreateTanggapan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idLaporan := vars["idLaporan"]
	var laporan models.Laporan

	if err := json.NewDecoder(r.Body).Decode(&laporan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if laporan.Tanggapan == "" || laporan.Status == "" {
		http.Error(w, "Gagal! Harap mengisi tanggapan", http.StatusBadRequest)
		return
	}

	laporan.TglTanggapan = time.Now().Format("2006-01-02")

	numRows := models.CreateTanggapan(idLaporan, laporan)
	if numRows == 0 {
		http.Error(w, "Gagal! Laporan tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}
