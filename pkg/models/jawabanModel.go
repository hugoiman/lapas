package models

import (
	"lapas/db"
)

// Jawaban is jawaban survei
type Jawaban struct {
	IDJawaban int `json:"idJawaban"`
	IDUser    int `json:"idUser"`
	IDSoal    int `json:"idSoal" validate:"required"`
	Jawaban   int `json:"jawaban" validate:"required,gte=1,lte=5"`
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
