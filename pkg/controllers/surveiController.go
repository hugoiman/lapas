package controllers

import (
	"encoding/json"
	models "lapas/pkg/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetSurvei is function
func GetSurvei(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	data, err := models.GetSurvei(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // survei not found
		return
	}

	message, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetSurveiActived is function
func GetSurveiActived(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	data, err := models.GetSurveiActived(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetSurveis is function
func GetSurveis(w http.ResponseWriter, r *http.Request) {
	data := models.GetSurveis()
	message, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateSurvei is new survei
func CreateSurvei(w http.ResponseWriter, r *http.Request) {
	var survei models.Survei
	if err := json.NewDecoder(r.Body).Decode(&survei); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err := validator.New().Struct(survei); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create slug
	var concateStr = survei.Judul + " " + survei.Periode.Format("January 2006")
	var replaceStr = strings.Replace(concateStr, " ", "-", -1)
	var slug = strings.ToLower(replaceStr)

	survei.Slug = slug
	idSurvei, err := models.CreateSurvei(survei)
	if err != nil {
		http.Error(w, "Gagal! Judul dan periode yang sama pernah dibuat sebelumnya.", http.StatusInternalServerError)
		return
	}

	survei.IDSurvei = idSurvei
	for _, v := range survei.Soal {
		if err = models.CreateSoal(idSurvei, v); err != nil {
			models.DeleteSurvei(strconv.Itoa(idSurvei))
			http.Error(w, "Gagal! Sub survei tidak terdaftar.", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// DeleteSurvei is delete survei
func DeleteSurvei(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]

	numRows := models.DeleteSurvei(idSurvei)
	if numRows == 0 {
		http.Error(w, "Gagal! Survei tidak ditemukan.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Survei berhasil dihapus!"}`))
}

// UpdateSurvei is Edit Survei
func UpdateSurvei(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]

	var survei models.Survei
	if err := json.NewDecoder(r.Body).Decode(&survei); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err := validator.New().Struct(survei); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create slug
	var concateStr = survei.Judul + " " + survei.Periode.Format("January 2006")
	var replaceStr = strings.Replace(concateStr, " ", "-", -1)
	var slug = strings.ToLower(replaceStr)

	survei.Slug = slug
	err := models.UpdateSurvei(idSurvei, survei)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // duplicate slug or survei not found
		return
	}

	for _, v := range survei.Soal {
		if err = models.UpdateSoal(idSurvei, v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // soal not found/not same or idsurvei not found/not same
			return
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// DuplicateSurvei is duplicate survei
func DuplicateSurvei(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]

	var survei models.Survei
	if err := json.NewDecoder(r.Body).Decode(&survei); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err := validator.New().Struct(survei); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if _, err := models.GetSurvei(idSurvei); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create slug
	var concateStr = survei.Judul + " " + survei.Periode.Format("January 2006")
	var replaceStr = strings.Replace(concateStr, " ", "-", -1)
	var slug = strings.ToLower(replaceStr)

	survei.Slug = slug
	idSurveiNew, err := models.CreateSurvei(survei)
	if err != nil {
		http.Error(w, "Gagal! Judul dan periode yang sama pernah dibuat sebelumnya.", http.StatusInternalServerError)
		return
	}

	// Get soal
	soals := models.GetSoal(idSurvei)

	survei.Soal = soals

	for _, v := range survei.Soal {
		_ = models.CreateSoal(idSurveiNew, v)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Survei berhasil diduplikasi!"}`))

}
