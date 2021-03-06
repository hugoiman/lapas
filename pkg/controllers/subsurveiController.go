package controllers

import (
	"encoding/json"
	models "lapas/pkg/models"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetSubSurvei is get list sub
func GetSubSurvei(w http.ResponseWriter, r *http.Request) {
	data := models.GetSubSurvei()
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateSubSurvei is new sub
func CreateSubSurvei(w http.ResponseWriter, r *http.Request) {
	var sub models.SubSurvei
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	models.CreateSubSurvei(sub)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// DeleteSubSurvei is update status sub
func DeleteSubSurvei(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSub := vars["idSub"]

	isDeleted := models.DeleteSubSurvei(idSub)
	if isDeleted == false {
		http.Error(w, "Gagal! Survei tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil dihapus!"}`))
}
