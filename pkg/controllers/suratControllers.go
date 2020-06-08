package controllers

import (
	"encoding/json"
	models "lapas/pkg/models"
	"net/http"
	"regexp"
	"strconv"
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
	surat.InputBy = strconv.Itoa(user.IDUser)

	if surat.Asal == "Internal" || surat.Tujuan == "Internal" {
		// next steps
	} else {
		http.Error(w, "Gagal! Asal/tujuan surat salah satunya harus berisi Internal", http.StatusBadRequest)
		return
	}

	penerima := models.GetUser(surat.Penerima)

	if penerima.Nama == "" {
		http.Error(w, "Gagal! User tidak ditemukan.", http.StatusBadRequest)
		return
	} else if penerima.Job != "Direksi" { // atau sekretaris perusahaan
		http.Error(w, "Gagal! Penerima tidak dizinkian.", http.StatusBadRequest)
		return
	} else if !penerima.Actived {
		http.Error(w, "Gagal! Penerima tidak aktif.", http.StatusBadRequest)
		return
	}

	// SendEmail(user.Email)

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

	surat.UpdatedAt = time.Now().Format("2006-01-02")
	surat.UpdatedBy = strconv.Itoa(user.IDUser)
	regexDate := regexp.MustCompile(`^(20)\d\d[-](0?[1-9]|1[012])[-](0?[1-9]|[12][0-9]|3[01])$`)

	if !regexDate.MatchString(surat.TglSurat) {
		http.Error(w, "Gagal! Format tanggal YYYY-MM-DD", http.StatusBadRequest)
		return
	} else if surat.TglDiterima != "" && !regexDate.MatchString(surat.TglDiterima) {
		http.Error(w, "Gagal! Format tanggal diterima YYYY-MM-DD", http.StatusBadRequest)
		return
	} else if surat.Status == "Waiting" {
		http.Error(w, "Gagal! Status surat harus 'Waiting'", http.StatusBadRequest)
		return
	}

	if surat.Asal == "Internal" || surat.Tujuan == "Internal" {
		// next steps
	} else {
		http.Error(w, "Gagal! Asal/tujuan surat salah satunya harus berisi Internal", http.StatusBadRequest)
		return
	}

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

	// SendEmail(user.Email)
	getSurat, _ := models.GetSurat(idSurat)
	if getSurat.Nomor == "" {
		http.Error(w, "Gagal! Surat tidak ditemukan.", http.StatusBadRequest)
		return
	} else if getSurat.Status == "Deleted" {
		http.Error(w, "Gagal! Surat telah dihapus.", http.StatusGone)
		return
	} else if getSurat.Status != "Waiting" {
		http.Error(w, "Gagal! Surat sudah ditindaklanjuti.", http.StatusBadRequest)
		return
	}

	err := models.UpdateSurat(idSurat, surat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // nomor tidak unik
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
	deletedBy := strconv.Itoa(user.IDUser)
	updatedAt := time.Now().Format("2006-01-02")

	surat, err := models.GetSurat(idSurat)

	if err != nil {
		http.Error(w, "Gagal! Surat tidak ditemukan.", http.StatusBadRequest)
		return
	} else if surat.Status != "Waiting" {
		http.Error(w, "Gagal! urat sudah ditindaklanjuti.", http.StatusBadRequest)
		return
	}

	models.DeleteSurat(idSurat, deletedBy, updatedAt)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil dihapus."}`))
}

// BeriStatus is func
func BeriStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurat := vars["idSurat"]
	user := context.Get(r, "user").(*MyClaims)

	surat, err := models.GetSurat(idSurat)
	if err != nil {
		http.Error(w, "Gagal! Surat tidak ditemukan.", http.StatusBadRequest)
		return
	} else if surat.Status == "Deleted" {
		http.Error(w, "Gagal! Surat telah dihapus.", http.StatusForbidden)
		return
	} else if surat.Penerima != strconv.Itoa(user.IDUser) {
		http.Error(w, "Gagal! Anda bukan penerima surat.", http.StatusForbidden)
		return
	}

	models.BeriStatus(idSurat)

	// filling
	// update status disposisi = delete, laporanDisposisi = delete where idsurat = x

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Status surat telah menjadi 'Filling'."}`))

}
