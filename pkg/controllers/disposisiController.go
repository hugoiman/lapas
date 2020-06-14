package controllers

import (
	"encoding/json"
	"fmt"
	models "lapas/pkg/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetDisposisi is func
func GetDisposisi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idDisposisi := vars["idDisposisi"]

	data, err := models.GetDisposisi(idDisposisi)
	if err != nil {
		http.Error(w, "Gagal! Surat tidak ditemukan", http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetDisposisis is func
func GetDisposisis(w http.ResponseWriter, r *http.Request) {
	data := models.GetDisposisis()
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetMyDisposisis is func
func GetMyDisposisis(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	data := models.GetMyDisposisis(strconv.Itoa(user.IDUser))
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// BeriStatusDisposisi is func
func BeriStatusDisposisi(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idDisposisi := vars["idDisposisi"]

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	if err := validator.New().Var(fmt.Sprintf("%v", data["status"]), "required,eq=Solved|eq=Waiting"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := data["status"].(string)

	disposisi, err := models.GetDisposisi(idDisposisi)
	surat, err := models.GetSurat(strconv.Itoa(disposisi.IDSurat))
	if err != nil {
		http.Error(w, "Gagal! Disposisi tidak ditemukan", http.StatusBadRequest)
		return
	} else if disposisi.IDPemberi != user.IDUser {
		http.Error(w, "Gagal! Anda bukan pembuat disposisi ini", http.StatusBadRequest)
		return
	} else if status == "Solved" && surat.IDPenerima != user.IDUser {
		http.Error(w, "Gagal! Hanya penerima surat yang dapat mengubah status menjadi 'Solved'.", http.StatusBadRequest)
		return
	}

	updatedByID := strconv.Itoa(user.IDUser)
	updatedAt := time.Now().Format("2006-01-02")
	models.BeriStatusDisposisi(idDisposisi, updatedByID, updatedAt, status)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Status menjadi 'Solved'!"}`))
}

// CreateDisposisi is func
func CreateDisposisi(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var disposisi models.Disposisi

	if err := json.NewDecoder(r.Body).Decode(&disposisi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(disposisi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var level = map[string]int{
		"Asistant Vice President": 6,
		"Manager":                 5,
		"Assistant Manager":       4,
		"Supervisor":              3,
		"Sr. Staff":               2,
		"Staff":                   1,
	}

	// list penerima
	penerima := []int{}
	for _, v := range disposisi.LaporanDispos {
		penerima = append(penerima, v.IDPenerima)
	}

	// IDPenerima = duplicate?
	keys := map[int]bool{}
	for v := range penerima {
		if x := keys[penerima[v]]; x {
			// Penerima dupplicate
			http.Error(w, "Gagal! Terdapat penerima yang sama.", http.StatusBadRequest)
			return
		} else if penerima[v] == user.IDUser {
			http.Error(w, "Gagal! Tidak dapat mendisposisikan ke diri sendiri.", http.StatusBadRequest)
			return
		}

		// Penerima tidak duplicate
		keys[penerima[v]] = true
	}

	for _, v := range disposisi.LaporanDispos {
		penerima := models.GetUser(strconv.Itoa(v.IDPenerima))
		if penerima.Nama == "" {
			http.Error(w, "Gagal! Penerima tidak ditemukan.", http.StatusBadRequest)
			return
		} else if !penerima.Actived {
			http.Error(w, "Gagal! Status penerima ("+penerima.Nama+") tidak aktif.", http.StatusBadRequest)
			return
		} else if (user.Divisi != "Sekretaris Perusahaan" && user.Job != "Direksi" && user.Job != "Direktur" && level[user.Pangkat] <= level[penerima.Pangkat]) || penerima.Job == "Direksi" || penerima.Job == "Direktur" {
			http.Error(w, "Gagal! Tidak dapat mendisposisikan ke pangkat yg lebih tinggi atau setara.", http.StatusBadRequest)
			return
		}
	}

	// jika bukan direksi -> cek apakah user pernah mendapat dan membuat disposisi pada surat tsb
	// cek apakah user pernah membuat disposisi pada surat tsb
	// jika direksi/sekretaris -> cek apakah surat sudah pernah didisposisikan

	surat, err := models.GetSurat(strconv.Itoa(disposisi.IDSurat))
	if err != nil {
		http.Error(w, "Gagal! Surat tidak ditemukan.", http.StatusBadRequest)
		return
	} else if len(surat.Disposisis) == 0 && surat.IDPenerima != user.IDUser {
		http.Error(w, "Gagal! Anda bukan penerima surat ini.", http.StatusForbidden)
		return
	}

	receiver := false
	for _, v := range surat.Disposisis {
		if user.IDUser == v.IDPemberi {
			http.Error(w, "Gagal! Anda sudah membuat disposisi pada surat ini.", http.StatusBadRequest)
			return
		}

		// jika user adalah penerima dispo
		for _, val := range v.LaporanDispos {
			if user.IDUser == val.IDPenerima {
				receiver = true
				break
			}
		}

		// user bukan penerima dispo sebelumnya? (u !direksi/direktur)
		if !receiver && len(surat.Disposisis) != 0 {
			http.Error(w, "Gagal! Anda bukan penerima disposisi pada surat ini.", http.StatusForbidden)
			return
		}
	}

	disposisi.IDPemberi = user.IDUser
	disposisi.Status = "Waiting"
	disposisi.CreatedAt = time.Now().Format("2006-01-02")

	idDisposisi := models.CreateDisposisi(disposisi)
	for _, v := range disposisi.LaporanDispos {
		v.IDDisposisi = idDisposisi
		v.Status = "Waiting"
		models.InitialLaporanDispo(user.IDUser, v)
	}

	if len(surat.Disposisis) == 0 {
		models.BeriStatusSurat(strconv.Itoa(disposisi.IDSurat), strconv.Itoa(user.IDUser), time.Now().Format("2006-01-02"), "Dispo")
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Berhasil menyimpan data disposisi."}`))
}
