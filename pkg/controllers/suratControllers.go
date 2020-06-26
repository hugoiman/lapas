package controllers

import (
	"encoding/json"
	"fmt"
	models "lapas/pkg/models"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetSurat is func
func GetSurat(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idSurat := vars["idSurat"]

	data, err := models.GetSurat(idSurat)
	if err != nil {
		http.Error(w, "Gagal! Surat tidak ditemukan", http.StatusBadRequest)
		return
	} else if user.Job != "Direktur" && user.Job != "Direksi" && user.Divisi != "Sekretaris Perusahaan" && user.Divisi != "Logistik & Administrasi" && user.IDUser != data.IDPenerima {
		http.Error(w, "Gagal! Anda tidak diizinkan.", http.StatusForbidden)
		return
	} else if data.Status == "Deleted" {
		http.Error(w, "Gagal! Surat telah dihapus.", http.StatusGone)
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
	regexDate := regexp.MustCompile(`^(20)\d\d[-](0?[1-9]|1[012])[-](0?[1-9]|[12][0-9]|3[01])$`)

	if err := json.NewDecoder(r.Body).Decode(&surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if surat.Asal != "Internal" && surat.Tujuan != "Internal" {
		http.Error(w, "Gagal! Asal/tujuan surat salah satunya harus berisi Internal", http.StatusBadRequest)
		return
	} else if !regexDate.MatchString(surat.TglSurat) {
		http.Error(w, "Gagal! Format tanggal surat harus YYYY-MM-DD", http.StatusBadRequest)
		return
	} else if surat.TglDiterima != "" && !regexDate.MatchString(surat.TglDiterima) {
		http.Error(w, "Gagal! Format tanggal diterima harus YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	surat.CreatedAt = time.Now().Format("2006-01-02")
	surat.Status = "Waiting"
	surat.InputByID = user.IDUser

	penerima := models.GetUser(strconv.Itoa(surat.IDPenerima))
	if penerima.Nama == "" {
		http.Error(w, "Gagal! Penerima tidak terdaftar.", http.StatusBadRequest)
		return
	} else if !(penerima.Job == "Direksi" || penerima.Job == "Direktur") && penerima.Divisi != "Sekretaris Perusahaan" {
		http.Error(w, "Gagal! Penerima tidak dizinkian.", http.StatusBadRequest)
		return
	} else if !penerima.Actived {
		http.Error(w, "Gagal! Penerima tidak aktif.", http.StatusBadRequest)
		return
	}

	err := models.CreateSurat(surat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // duplikasi nomor/lampiran (not unique)
		return
	}

	subject := "Surat Masuk"
	message := "<p><b>[New]</b> - Surat dari " + surat.Asal + " ke " + surat.Tujuan +
		".<br>No: " + surat.Nomor +
		".<br>Tanggal Surat: " + surat.TglSurat +
		".<br>Perihal: " + surat.Perihal + "</p>"
	address := []string{penerima.Email}
	SendEmail(subject, address, message)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Data berhasil disimpan!"}`))
}

// UpdateSurat dilakukan ketika status == waiting
func UpdateSurat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurat := vars["idSurat"]
	user := context.Get(r, "user").(*MyClaims)
	var surat models.Surat
	regexDate := regexp.MustCompile(`^(20)\d\d[-](0?[1-9]|1[012])[-](0?[1-9]|[12][0-9]|3[01])$`)

	if err := json.NewDecoder(r.Body).Decode(&surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(surat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if surat.Asal != "Internal" && surat.Tujuan != "Internal" {
		http.Error(w, "Gagal! Asal/tujuan surat salah satunya harus berisi Internal", http.StatusBadRequest)
		return
	} else if !regexDate.MatchString(surat.TglSurat) {
		http.Error(w, "Gagal! Format tanggal surat harus YYYY-MM-DD", http.StatusBadRequest)
		return
	} else if surat.TglDiterima != "" && !regexDate.MatchString(surat.TglDiterima) {
		http.Error(w, "Gagal! Format tanggal diterima harus YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	surat.UpdatedAt = time.Now().Format("2006-01-02")
	surat.UpdatedByID = user.IDUser

	penerima := models.GetUser(strconv.Itoa(surat.IDPenerima))
	if penerima.Nama == "" {
		http.Error(w, "Gagal! Penerima tidak terdaftar.", http.StatusBadRequest)
		return
	} else if !(penerima.Job == "Direksi" || penerima.Job == "Direktur") && penerima.Divisi != "Sekretaris Perusahaan" {
		http.Error(w, "Gagal! Penerima tidak dizinkian.", http.StatusBadRequest)
		return
	} else if !penerima.Actived {
		http.Error(w, "Gagal! Penerima tidak aktif.", http.StatusBadRequest)
		return
	}

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
		http.Error(w, err.Error(), http.StatusBadRequest) // duplikasi nomor/lampiran (not unique)
		return
	}

	subject := "Surat Masuk"
	message := "<p><b>[Updated]</b> - Surat dari " + surat.Asal + " ke " + surat.Tujuan +
		".<br>No: " + surat.Nomor +
		".<br>Tanggal Surat: " + surat.TglSurat +
		".<br>Perihal: " + surat.Perihal + "</p>"
	address := []string{penerima.Email}
	SendEmail(subject, address, message)

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
		http.Error(w, "Gagal! Surat sudah ditindaklanjuti.", http.StatusBadRequest)
		return
	}

	models.DeleteSurat(idSurat, deletedBy, updatedAt)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil dihapus."}`))
}

// BeriStatusSurat is func
func BeriStatusSurat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurat := vars["idSurat"]
	user := context.Get(r, "user").(*MyClaims)

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	if err := validator.New().Var(fmt.Sprintf("%v", data["status"]), "required,eq=Undelete|eq=Filling"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := data["status"].(string)

	surat, err := models.GetSurat(idSurat)
	if err != nil {
		http.Error(w, "Gagal! Surat tidak ditemukan.", http.StatusBadRequest)
		return
	} else if status == "Undelete" && user.Job == "Direksi" || user.Job == "Direktur" {
		http.Error(w, "Gagal! Anda tidak diizinkan.", http.StatusForbidden)
		return
	} else if status == "Undelete" && surat.Status != "Deleted" {
		http.Error(w, "Gagal! Surat belum dihapus sebelumnya.", http.StatusBadRequest)
		return
	} else if status == "Filling" && surat.IDPenerima != user.IDUser {
		http.Error(w, "Gagal! Anda bukan penerima surat.", http.StatusBadRequest)
		return
	} else if surat.Status == "Solved" {
		http.Error(w, "Gagal! Status surat sudah 'Solved'.", http.StatusBadRequest)
		return
	}

	updatedByID := user.IDUser
	updatedAt := time.Now().Format("2006-01-02")

	models.BeriStatusSurat(idSurat, strconv.Itoa(updatedByID), updatedAt, status)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Status surat telah menjadi 'Filling'."}`))

}
