package controllers

import (
	"encoding/json"
	"fmt"
	models "lapas/pkg/models"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// SaveJawaban is simpan jawaban
func SaveJawaban(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]
	idUser := vars["idUser"]

	var jawaban models.Jawabans
	if err := json.NewDecoder(r.Body).Decode(&jawaban); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err := validator.New().Struct(jawaban); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idSoal := make([]int, 0)
	for _, v := range jawaban.Jawabans {
		idSoal = append(idSoal, v.IDSoal)
	}

	totalID := len(idSoal)

	// idSoal duplikat?
	for i := 0; i < totalID; i++ {
		for j := i + 1; j < totalID; j++ {
			if idSoal[i] == idSoal[j] {
				http.Error(w, "Gagal! Ditemukan duplikasi soal!", http.StatusInternalServerError)
				return
			}
		}
	}

	// survei ditemukan?
	survei, err := models.GetSurvei(idSurvei)
	if err != nil {
		http.Error(w, "Gagal! Survei tidak ditemukan!", http.StatusInternalServerError)
		return
	}

	// jumlah jawaban sama dengan soal?
	totalSoal := len(survei.Soal)
	if totalID != totalSoal {
		http.Error(w, "Gagal! Jumlah jawaban berbeda dari soal survei!", http.StatusInternalServerError)
		return
	}

	// id soal json sama dengan id soal pada survei?
	idSoalSurvei := make([]int, 0)
	for _, v := range survei.Soal {
		idSoalSurvei = append(idSoalSurvei, v.IDSoal)
	}

	sameID := false
loop:
	for i := 1; i < totalID; i++ {
		for j := 1; j < totalID; j++ {
			if idSoal[i] == idSoalSurvei[j] {
				sameID = true
				continue loop
			} else {
				sameID = false
			}
		}
	}

	if sameID == false {
		http.Error(w, "Gagal! Terdapat soal yang berbeda dari soal survei!", http.StatusInternalServerError)
		return
	}

	// survei sudah pernah diisi?
	var message string
	getJawaban := models.GetJawaban(idSurvei, idUser)
	if len(getJawaban.Jawabans) == 0 {
		for _, v := range jawaban.Jawabans {
			fmt.Printf("%+v", v)
			// models.CreateJawaban(idUser, v)
		}
		message = `{"message":"Simpan!"}`
	} else {
		success, message := UpdateJawaban(idUser, idSurvei, jawaban, getJawaban)
		// models.UpdateJawaban(jawaban)
		if success == false {
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

// UpdateJawaban is edit jawaban
func UpdateJawaban(idUser, idSurvei string, newJawaban models.Jawabans, oldJawaban models.Jawabans) (bool, string) {
	idNewJawaban := make([]int, 0)
	for _, v := range newJawaban.Jawabans {
		idNewJawaban = append(idNewJawaban, v.IDJawaban)
	}

	totalID := len(idNewJawaban)

	// idJawaban duplikat?
	for i := 0; i < totalID; i++ {
		for j := i + 1; j < totalID; j++ {
			if idNewJawaban[i] == idNewJawaban[j] {
				message := "Gagal! Ditemukan duplikasi jawaban!"
				return false, message
			}
		}
	}

	// id jawaban json sudah sesuai dgn id jawaban di table?
	idOldJawaban := make([]int, 0)
	for _, v := range oldJawaban.Jawabans {
		idOldJawaban = append(idOldJawaban, v.IDJawaban)
	}

	sameID := false
loop:
	for i := 1; i < totalID; i++ {
		for j := 1; j < totalID; j++ {
			if idNewJawaban[i] == idOldJawaban[j] {
				sameID = true
				continue loop
			} else {
				sameID = false
			}
		}
	}

	if sameID == false {
		message := "Gagal! Terdapt ID jawaban tidak sesuai dengan yang ada!"
		return false, message
	}

	for _, v := range newJawaban.Jawabans {
		fmt.Printf("%+v", v)
		// models.UpdateJawaban(v)
	}

	return true, "Jawaban berhasil di simpan"
}
