package controllers

import (
	"encoding/json"
	"fmt"
	models "lapas/pkg/models"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetSurat is func
func GetSurat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurat := vars["idSurat"]

	data, err := models.GetSurat(idSurat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if strings.EqualFold(data.Status, "Deleted") {
		http.Error(w, "Surat telah dihapus.", http.StatusGone)
		return
	}
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetSurats is get list surat
func GetSurats(w http.ResponseWriter, r *http.Request) {
	data := models.GetSurats()
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateSurat is func
func CreateSurat(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var surat models.Surat

	if err := json.NewDecoder(r.Body).Decode(&surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regexDate := regexp.MustCompile(`^(20)\d\d[-](0?[1-9]|1[012])[-](0?[1-9]|[12][0-9]|3[01])$`)

	if !regexDate.MatchString(surat.TglSurat) {
		http.Error(w, "Gagal! Format tanggal YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if surat.TglDiterima != "" && !regexDate.MatchString(surat.TglDiterima) {
		http.Error(w, "Gagal! Format tanggal trima YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	surat.CreatedAt = time.Now().Format("2006-01-02")
	surat.Status = "Waiting"
	surat.InputBy = user.Nama

	if surat.Asal == "Internal" || surat.Tujuan == "Internal" {
		// next steps
	} else {
		http.Error(w, "Gagal! Asal/tujuan surat salah satunya harus berisi Internal", http.StatusBadRequest)
		return
	}

	if strings.EqualFold(surat.Penerima, "Pimpinan") {
		// get email pimpinan
		penerima := models.GetPimpinan()

		for _, v := range penerima.Users {
			fmt.Println(v.Nama, " : ", v.Email)
			// SendEmail(v.Email)
		}

	} else {
		penerima := models.GetUser(surat.Penerima)

		if penerima.Nama == "" {
			http.Error(w, "Gagal! User tidak ditemukan.", http.StatusBadRequest)
			return
		} else if penerima.Job != "Direksi" {
			http.Error(w, "Gagal! Penerima tidak dizinkian.", http.StatusBadRequest)
			return
		} else if !penerima.Actived {
			http.Error(w, "Gagal! Penerima tidak aktif.", http.StatusBadRequest)
			return
		}

		surat.Penerima = penerima.Nama

		// SendEmail(user.Email)

	}

	err := models.CreateSurat(surat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// UpdateSurat is func
func UpdateSurat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurat := vars["idSurat"]
	user := context.Get(r, "user").(*MyClaims)
	var surat models.Surat

	if err := json.NewDecoder(r.Body).Decode(&surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regexDate := regexp.MustCompile(`^(20)\d\d[-](0?[1-9]|1[012])[-](0?[1-9]|[12][0-9]|3[01])$`)

	if !regexDate.MatchString(surat.TglSurat) {
		http.Error(w, "Gagal! Format tanggal YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if surat.TglDiterima != "" && !regexDate.MatchString(surat.TglDiterima) {
		http.Error(w, "Gagal! Format tanggal trima YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	surat.UpdatedAt = time.Now().Format("2006-01-02")
	surat.UpdatedBy = user.Nama

	if surat.Asal == "Internal" || surat.Tujuan == "Internal" {
		// next steps
	} else {
		http.Error(w, "Gagal! Asal/tujuan surat salah satunya harus berisi Internal", http.StatusBadRequest)
		return
	}

	if strings.EqualFold(surat.Penerima, "Pimpinan") {
		// get email pimpinan
		penerima := models.GetPimpinan()

		for _, v := range penerima.Users {
			fmt.Println(v.Nama, " : ", v.Email)
			// SendEmail(v.Email)
		}

	} else {
		penerima := models.GetUser(surat.Penerima)

		if penerima.Nama == "" {
			http.Error(w, "Gagal! User tidak ditemukan.", http.StatusBadRequest)
			return
		} else if penerima.Job != "Direksi" {
			http.Error(w, "Gagal! Penerima tidak dizinkian.", http.StatusBadRequest)
			return
		} else if !penerima.Actived {
			http.Error(w, "Gagal! Penerima tidak aktif.", http.StatusBadRequest)
			return
		}

		surat.Penerima = penerima.Nama

		// SendEmail(user.Email)

	}

	numRows := models.UpdateSurat(idSurat, surat)

	if numRows == 0 {
		http.Error(w, "Gagal! Surat tidak ditemukan atau sudah dihapus.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil diperbarui."}`))
}

// DeleteSurat is func
func DeleteSurat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurat := vars["idSurat"]
	user := context.Get(r, "user").(*MyClaims)
	deletedBy := user.Nama
	updatedAt := time.Now().Format("2006-01-02")

	numRows := models.DeleteSurat(idSurat, deletedBy, updatedAt)
	if numRows == 0 {
		http.Error(w, "Gagal! Surat tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil dihapus."}`))
}
