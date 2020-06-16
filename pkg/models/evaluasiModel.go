package models

import (
	"lapas/db"
	"time"
)

// Evaluasi is class
type Evaluasi struct {
	IDEvaluasi           int       `json:"idEvaluasi"`
	IDSurvei             int       `json:"idSurvei" validate:"required"`
	Pesan                string    `json:"pesan"`
	Utama                string    `json:"utama"`
	TeknikPengembangan   string    `json:"teknik_pengembangan"`
	Operasi              string    `json:"operasi"`
	KeuanganAdministrasi string    `json:"keuangan_administrasi"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// GetEvaluasi is func
func GetEvaluasi(idSurvei string) (Evaluasi, error) {
	con := db.Connect()
	query := "SELECT idEvaluasi, idSurvei, pesan, utama, teknik_pengembangan, operasi, keuangan_administrasi, createdAt, updatedAt FROM evaluasi WHERE idSurvei = ?"

	evaluasi := Evaluasi{}
	err := con.QueryRow(query, idSurvei).Scan(
		&evaluasi.IDEvaluasi, &evaluasi.IDSurvei, &evaluasi.Pesan, &evaluasi.Utama, &evaluasi.TeknikPengembangan,
		&evaluasi.Operasi, &evaluasi.KeuanganAdministrasi, &evaluasi.CreatedAt, &evaluasi.UpdatedAt)

	defer con.Close()
	return evaluasi, err
}

// CreateEvaluasi is new Evaluasi
func CreateEvaluasi(evaluasi Evaluasi) error {
	con := db.Connect()
	_, err := con.Exec("INSERT INTO evaluasi (idSurvei, pesan, utama, teknik_pengembangan, operasi, keuangan_administrasi, createdAt) VALUES (?,?,?,?,?,?,?)", evaluasi.IDSurvei, evaluasi.Pesan, &evaluasi.Utama, &evaluasi.TeknikPengembangan, &evaluasi.Operasi, &evaluasi.KeuanganAdministrasi, evaluasi.CreatedAt)

	defer con.Close()

	return err
}

// UpdateEvaluasi is edit Evaluasi
func UpdateEvaluasi(idEvaluasi string, evaluasi Evaluasi) bool {
	con := db.Connect()
	query := "UPDATE evaluasi SET pesan = ?, utama = ?, teknik_pengembangan = ?, operasi = ?, keuangan_administrasi = ?, updatedAt = ? WHERE idEvaluasi = ? AND idSurvei = ?"
	res, _ := con.Exec(query, evaluasi.Pesan, &evaluasi.Utama, &evaluasi.TeknikPengembangan, &evaluasi.Operasi, &evaluasi.KeuanganAdministrasi, evaluasi.UpdatedAt, idEvaluasi, evaluasi.IDSurvei)

	count, _ := res.RowsAffected()

	defer con.Close()

	if int(count) == 0 {
		return false
	}

	return true
}
