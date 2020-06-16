package controllers

import (
	"encoding/json"
	models "lapas/pkg/models"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetSurvei is function
func GetSurvei(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	data, err := models.GetSurvei(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // survei not found
		return
	}

	message, _ := json.Marshal(data)

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetSurveis is function
func GetSurveis(w http.ResponseWriter, r *http.Request) {
	data := models.GetSurveis()
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateSurvei is new survei
func CreateSurvei(w http.ResponseWriter, r *http.Request) {
	var survei models.Survei
	if err := json.NewDecoder(r.Body).Decode(&survei); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(survei); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create slug
	var concateStr = survei.Judul + " " + survei.Periode.Format("January 2006")
	var replaceStr = strings.Replace(concateStr, " ", "-", -1)
	var slug = strings.ToLower(replaceStr)

	survei.Slug = slug
	idSurvei, err := models.CreateSurvei(survei)
	if err != nil {
		http.Error(w, "Gagal! Judul dan periode yang sama pernah dibuat sebelumnya.", http.StatusBadRequest)
		return
	}

	survei.IDSurvei = idSurvei
	for _, v := range survei.Soal {
		if err = models.CreateSoal(idSurvei, v); err != nil {
			_ = models.DeleteSurvei(strconv.Itoa(idSurvei))
			http.Error(w, "Gagal! Sub survei tidak terdaftar.", http.StatusBadRequest)
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

	isDeleted := models.DeleteSurvei(idSurvei)
	if isDeleted == false {
		http.Error(w, "Gagal! Survei tidak ditemukan.", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(survei); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create slug
	var concateStr = survei.Judul + " " + survei.Periode.Format("January 2006")
	var replaceStr = strings.Replace(concateStr, " ", "-", -1)
	var slug = strings.ToLower(replaceStr)

	survei.Slug = slug
	err := models.UpdateSurvei(idSurvei, survei)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // duplicate slug or survei not found
		return
	}

	for _, v := range survei.Soal {
		if err = models.UpdateSoal(idSurvei, v); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // soal not found/not same or idsurvei not found/not same
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(survei); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, err := models.GetSurvei(idSurvei); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create slug
	var concateStr = survei.Judul + " " + survei.Periode.Format("January 2006")
	var replaceStr = strings.Replace(concateStr, " ", "-", -1)
	var slug = strings.ToLower(replaceStr)

	survei.Slug = slug
	idSurveiNew, err := models.CreateSurvei(survei)
	if err != nil {
		http.Error(w, "Gagal! Judul dan periode yang sama pernah dibuat sebelumnya.", http.StatusBadRequest)
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

// ChangeStatus is change status survei
func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]

	survei, err := models.GetSurvei(idSurvei)
	if err != nil {
		http.Error(w, "Gagal! Survei tidak ditemukan.", http.StatusBadRequest) // survei not found
		return
	}

	actived := !survei.Actived
	strActived := strconv.FormatBool(actived)

	models.ChangeStatus(idSurvei, actived)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Status survei menjadi ` + strActived + ` !"}`))

}

// GetStatistikResponden is func
func GetStatistikResponden(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]

	data := models.GetDataResponden(idSurvei)
	var utama, teknikPengembangan, operasi, keuanganAdministrasi int
	var staff, supervisor, assManager, manager, avp int
	var u25, u26, u31, u36, u41 int

	// get direktorat
	for _, v := range data.Users {
		if strings.EqualFold(v.Direktorat, "utama") {
			utama++
		} else if strings.EqualFold(v.Direktorat, "teknik dan pengembangan") {
			teknikPengembangan++
		} else if strings.EqualFold(v.Direktorat, "operasi") {
			operasi++
		} else if strings.EqualFold(v.Direktorat, "keuangan dan administrasi") {
			keuanganAdministrasi++
		}
	}

	// get pangkat
	for _, v := range data.Users {
		if strings.EqualFold(v.Pangkat, "staff") || strings.EqualFold(v.Pangkat, "sr. staff") {
			staff++
		} else if strings.EqualFold(v.Pangkat, "supervisor") {
			supervisor++
		} else if strings.EqualFold(v.Pangkat, "assistant manager") {
			assManager++
		} else if strings.EqualFold(v.Pangkat, "manager") {
			manager++
		} else if strings.EqualFold(v.Pangkat, "assistant vice president") {
			avp++
		}
	}

	// get usia
	today := time.Now()

	for _, v := range data.Users {
		usia := math.Floor(today.Sub(v.TglLahir).Hours() / 24 / 365)
		if usia <= 25 {
			u25++
		} else if usia <= 30 {
			u26++
		} else if usia <= 35 {
			u31++
		} else if usia <= 40 {
			u36++
		} else if usia >= 41 {
			u41++
		}
	}

	message := []byte(`{
				"responden": ` + strconv.Itoa(len(data.Users)) + `,
				"direktorat":[
					{
						"utama": ` + strconv.Itoa(utama) + `,
						"teknik_pengembangan": ` + strconv.Itoa(teknikPengembangan) + `,
						"operasi": ` + strconv.Itoa(operasi) + `,
						"keuangan_administrasi": ` + strconv.Itoa(keuanganAdministrasi) + `
					}
				],
				"pangkat":[
					{
						"staff": ` + strconv.Itoa(staff) + `,
						"supervisor": ` + strconv.Itoa(supervisor) + `,
						"assistant_manager": ` + strconv.Itoa(assManager) + `,
						"manager": ` + strconv.Itoa(manager) + `,
						"avp": ` + strconv.Itoa(avp) + `
					}
				],
				"usia":[
					{
						"u25": ` + strconv.Itoa(u25) + `,
						"u26": ` + strconv.Itoa(u26) + `,
						"u31": ` + strconv.Itoa(u31) + `,
						"u36": ` + strconv.Itoa(u36) + `,
						"u41": ` + strconv.Itoa(u41) + `
					}
				]
			}`)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// Statistik is struct
type Statistik struct {
	STS int `json:"sts"`
	TS  int `json:"ts"`
	N   int `json:"n"`
	S   int `json:"s"`
	SS  int `json:"ss"`
}

// GetStatistikJawaban is func
func GetStatistikJawaban(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]
	direktorat := vars["direktorat"]
	direktorat = strings.Replace(direktorat, "-", " ", -1)

	var listDirektorat = [5]string{"semua", "utama", "teknik dan pengembangan", "operasi", "keuangan dan administrasi"}
	var isValid = false

	for _, v := range listDirektorat {
		if strings.EqualFold(v, direktorat) {
			isValid = true
		}
	}

	if !isValid {
		http.Error(w, "Gagal! Direktorat tidak terdaftar.", http.StatusBadRequest)
		return
	}

	statistik := make(map[string]Statistik)

	data := models.GetJawabans(idSurvei, direktorat)

	for _, v := range data.Jawabans {
		if v.Jawaban == 1 {
			statistik[v.SubSurvei] = Statistik{
				STS: statistik[v.SubSurvei].STS + 1,
				TS:  statistik[v.SubSurvei].TS + 0,
				N:   statistik[v.SubSurvei].N + 0,
				S:   statistik[v.SubSurvei].S + 0,
				SS:  statistik[v.SubSurvei].SS + 0,
			}
		} else if v.Jawaban == 2 {
			statistik[v.SubSurvei] = Statistik{
				STS: statistik[v.SubSurvei].STS + 0,
				TS:  statistik[v.SubSurvei].TS + 1,
				N:   statistik[v.SubSurvei].N + 0,
				S:   statistik[v.SubSurvei].S + 0,
				SS:  statistik[v.SubSurvei].SS + 0,
			}
		} else if v.Jawaban == 3 {
			statistik[v.SubSurvei] = Statistik{
				STS: statistik[v.SubSurvei].STS + 0,
				TS:  statistik[v.SubSurvei].TS + 0,
				N:   statistik[v.SubSurvei].N + 1,
				S:   statistik[v.SubSurvei].S + 0,
				SS:  statistik[v.SubSurvei].SS + 0,
			}
		} else if v.Jawaban == 4 {
			statistik[v.SubSurvei] = Statistik{
				STS: statistik[v.SubSurvei].STS + 0,
				TS:  statistik[v.SubSurvei].TS + 0,
				N:   statistik[v.SubSurvei].N + 0,
				S:   statistik[v.SubSurvei].S + 1,
				SS:  statistik[v.SubSurvei].SS + 0,
			}
		} else if v.Jawaban == 5 {
			statistik[v.SubSurvei] = Statistik{
				STS: statistik[v.SubSurvei].STS + 0,
				TS:  statistik[v.SubSurvei].TS + 0,
				N:   statistik[v.SubSurvei].N + 0,
				S:   statistik[v.SubSurvei].S + 0,
				SS:  statistik[v.SubSurvei].SS + 1,
			}
		}
	}

	message, _ := json.Marshal(statistik)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)

}

// GetDataResponden is func
func GetDataResponden(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idSurvei := vars["idSurvei"]

	data := models.GetDataResponden(idSurvei)
	message, _ := json.Marshal(data)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
