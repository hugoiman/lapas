package models

import (
	"lapas/db"
	"time"
)

// LaporanDispo is class
type LaporanDispo struct {
	IDLD        int    `json:"idLD"`
	IDDisposisi int    `json:"idDisposisi"`
	IDPenerima  int    `json:"idPenerima"`
	Penerima    string `json:"penerima"`
	Status      string `json:"status"`
	Laporan     string `json:"laporan" validate:"required"`
	Lampiran    string `json:"lampiran"`
	UpdatedAt   string `json:"updatedAt"`
}

// LaporanDispos is LaporanDispo List
type LaporanDispos struct {
	LaporanDispos []LaporanDispo `json:"laporan_dispo"`
}

// GetLaporanDispo is func
func GetLaporanDispo(idLaporan string) (LaporanDispo, error) {
	con := db.Connect()
	query := "SELECT IDLD, idDisposisi, idPenerima, (SELECT nama FROM user WHERE idUser = idPenerima) AS penerima, status, laporan, lampiran, updatedAt FROM laporan_disposisi WHERE idLD = ?"
	laporan := LaporanDispo{}

	var updatedAt interface{}

	err := con.QueryRow(query, idLaporan).Scan(
		&laporan.IDLD, &laporan.IDDisposisi, &laporan.IDPenerima, &laporan.Penerima, &laporan.Status, &laporan.Laporan, &laporan.Lampiran, &updatedAt)

	if updatedAt == nil {
		laporan.UpdatedAt = ""
	} else {
		laporan.UpdatedAt = updatedAt.(time.Time).Format("02 Jan 2006")
	}

	defer con.Close()
	return laporan, err
}

// GetLaporanDispos is func
func GetLaporanDispos(idDisposisi string) LaporanDispos {
	con := db.Connect()
	query := "SELECT IDLD, idDisposisi, idPenerima, (SELECT nama FROM user WHERE idUser = idPenerima) AS penerima, status, laporan, lampiran, updatedAt FROM laporan_disposisi WHERE idDisposisi = ?"
	rows, _ := con.Query(query, idDisposisi)

	laporanDispo := LaporanDispo{}
	laporanDispos := LaporanDispos{}

	var updatedAt interface{}

	for rows.Next() {
		_ = rows.Scan(
			&laporanDispo.IDLD, &laporanDispo.IDDisposisi, &laporanDispo.IDPenerima, &laporanDispo.Penerima, &laporanDispo.Status,
			&laporanDispo.Laporan, &laporanDispo.Lampiran, &updatedAt)

		if updatedAt == nil {
			laporanDispo.UpdatedAt = ""
		} else {
			laporanDispo.UpdatedAt = updatedAt.(time.Time).Format("02 Jan 2006")
		}

		laporanDispos.LaporanDispos = append(laporanDispos.LaporanDispos, laporanDispo)
	}

	defer con.Close()
	return laporanDispos
}

// InitialLaporanDispo is func
func InitialLaporanDispo(idPemberi int, laporan LaporanDispo) {
	con := db.Connect()
	_, _ = con.Exec("INSERT INTO laporan_disposisi (idDisposisi, idPenerima, status) VALUES (?,?,?)", laporan.IDDisposisi, laporan.IDPenerima, laporan.Status)

	_, _ = con.Exec("UPDATE laporan_disposisi status = 'Forward' WHERE idDisposisi = ? AND idPenerima = ?", laporan.IDDisposisi, idPemberi)

	defer con.Close()
}

// CreateLaporanDisposisi is func
func CreateLaporanDisposisi(idLaporan string, laporan LaporanDispo) {
	con := db.Connect()
	query := "UPDATE laporan_disposisi SET status = ?, laporan = ?, lampiran = ?, updatedAt = ? WHERE idLD = ?"
	_, _ = con.Exec(query, laporan.Status, laporan.Laporan, laporan.Lampiran, laporan.UpdatedAt, laporan.IDLD)

	query = "UPDATE disposisi SET status = 'Process' WHERE idDisposisi = ? AND idLD = ?"
	_, _ = con.Exec(query, laporan.Status, laporan.IDDisposisi)

	defer con.Close()
}
