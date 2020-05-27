package controllers

import (
	"encoding/json"
	models "lapas/pkg/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetEvaluasi is func
func GetEvaluasi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]

	data, err := models.GetEvaluasi(idSurvei)
	if err != nil {
		http.Error(w, "Gagal! Evaluasi tidak ditemukan.", http.StatusInternalServerError)
		return
	}

	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateEvaluasi is add new evaluasi
func CreateEvaluasi(w http.ResponseWriter, r *http.Request) {
	var evaluasi models.Evaluasi
	if err := json.NewDecoder(r.Body).Decode(&evaluasi); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err := validator.New().Struct(evaluasi); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	evaluasi.CreatedAt = time.Now()

	err := models.CreateEvaluasi(evaluasi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// UpdateEvaluasi is edit evaluasi
func UpdateEvaluasi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idEvaluasi := vars["idEvaluasi"]

	var evaluasi models.Evaluasi
	if err := json.NewDecoder(r.Body).Decode(&evaluasi); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err := validator.New().Struct(evaluasi); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	evaluasi.UpdatedAt = time.Now()

	numRows := models.UpdateEvaluasi(idEvaluasi, evaluasi)
	if numRows == 0 {
		http.Error(w, "Gagal! survei atau id evaluasi tidak ditemukan.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}
