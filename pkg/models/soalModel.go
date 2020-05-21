package models

import "lapas/db"

// Soal is part of survei
type Soal struct {
	IDSoal    int    `json:"idSoal"`
	IDSub     int    `json:"idSub" validate:"required"`
	Soal      string `json:"soal" validate:"required"`
	Subsurvei string `json:"subSurvei"`
}

// GetSoal is get soal
func GetSoal(idSurvei string) []Soal {
	con := db.Connect()
	query := "SELECT a.idSoal, a.idSub, a.soal, b.subsurvei FROM soal a JOIN subsurvei b ON a.idSub = b.idSub WHERE a.idSurvei = ?"
	rows, _ := con.Query(query, idSurvei)

	soal := Soal{}
	soals := []Soal{}

	for rows.Next() {
		_ = rows.Scan(&soal.IDSoal, &soal.IDSub, &soal.Soal, &soal.Subsurvei)
		soals = append(soals, soal)
	}

	defer con.Close()

	return soals
}

// CreateSoal is add soal
func CreateSoal(idSurvei int, soal Soal) error {
	con := db.Connect()
	_, err := con.Exec("INSERT INTO soal (idSurvei, idSub, soal) VALUES (?,?,?)", idSurvei, soal.IDSub, soal.Soal)

	defer con.Close()

	return err
}

// UpdateSoal is edit soal
func UpdateSoal(idSurvei string, soal Soal) error {
	con := db.Connect()
	query := "UPDATE soal SET idSub = ?, soal = ? WHERE idSurvei = ? AND idSoal = ?"
	_, err := con.Exec(query, soal.IDSub, soal.Soal, idSurvei, soal.IDSoal)

	defer con.Close()

	return err
}
