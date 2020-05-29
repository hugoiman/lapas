package models

import (
	"lapas/db"
)

// Jawaban is jawaban survei
type Jawaban struct {
	IDJawaban int    `json:"idJawaban"`
	IDUser    int    `json:"idUser"`
	IDSoal    int    `json:"idSoal" validate:"required"`
	Jawaban   int    `json:"jawaban" validate:"required,gte=1,lte=5"`
	SubSurvei string `json:"subSurvei"`
	Nama      string `json:"nama"`
}

// Jawabans is jawaban list
type Jawabans struct {
	Jawabans []Jawaban `json:"jawaban" validate:"required,dive"`
}

// GetJawaban is simpan jawaban
func GetJawaban(idSurvei, idUser string) Jawabans {
	con := db.Connect()
	query := "SELECT a.idJawaban, a.idUser, a.idSoal, a.jawaban FROM jawaban a JOIN soal b ON a.idSoal = b.idSoal WHERE a.idUser = ? AND b.idSurvei = ?"
	rows, _ := con.Query(query, idUser, idSurvei)

	jawaban := Jawaban{}
	jawabans := Jawabans{}

	for rows.Next() {
		_ = rows.Scan(&jawaban.IDJawaban, &jawaban.IDUser, &jawaban.IDSoal, &jawaban.Jawaban)
		jawabans.Jawabans = append(jawabans.Jawabans, jawaban)
	}

	defer con.Close()

	return jawabans
}

// CreateJawaban is new save new jawaban
func CreateJawaban(idUser string, jawaban Jawaban) {
	con := db.Connect()
	_, _ = con.Exec("INSERT INTO jawaban (idUser, idSoal, jawaban) VALUES (?,?,?)", idUser, jawaban.IDSoal, jawaban.Jawaban)

	defer con.Close()
}

// UpdateJawaban is edit Jawaban
func UpdateJawaban(jawaban Jawaban) {
	con := db.Connect()
	query := "UPDATE jawaban SET jawaban = ? WHERE idJawaban = ?"
	_, _ = con.Exec(query, jawaban.Jawaban, jawaban.IDJawaban)

	defer con.Close()
}

// GetJawabans is func
func GetJawabans(idSurvei, direktorat string) Jawabans {
	con := db.Connect()
	queryDirektorat := "SELECT a.idJawaban, a.idSoal, a.jawaban, c.subSurvei, d.nama FROM jawaban a JOIN soal b ON a.idSoal = b.idSoal JOIN subsurvei c ON b.idSub = c.idSub JOIN user d ON a.idUser = d.idUser WHERE b.idSurvei = ? AND d.direktorat = ?"
	query := "SELECT a.idJawaban, a.idSoal, a.jawaban, c.subSurvei, d.nama FROM jawaban a JOIN soal b ON a.idSoal = b.idSoal JOIN subsurvei c ON b.idSub = c.idSub JOIN user d ON a.idUser = d.idUser WHERE b.idSurvei = ?"
	rows, _ := con.Query(query, idSurvei)

	if direktorat != "semua" {
		rows, _ = con.Query(queryDirektorat, idSurvei, direktorat)
	}

	jawaban := Jawaban{}
	jawabans := Jawabans{}

	for rows.Next() {
		_ = rows.Scan(
			&jawaban.IDJawaban, &jawaban.IDSoal, &jawaban.Jawaban, &jawaban.SubSurvei, &jawaban.Nama)
		jawabans.Jawabans = append(jawabans.Jawabans, jawaban)
	}

	defer con.Close()
	return jawabans
}
