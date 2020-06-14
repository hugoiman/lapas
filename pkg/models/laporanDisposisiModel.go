package models

import (
	"fmt"
	"lapas/db"
	"time"
)

// LaporanDispo is class
type LaporanDispo struct {
	IDLD        int    `json:"idLD"`
	IDDisposisi int    `json:"idDisposisi"`
	IDPenerima  int    `json:"idPenerima" validate:"required"`
	Penerima    string `json:"penerima"`
	Status      string `json:"status"`
	Laporan     string `json:"laporan"`
	Lampiran    string `json:"lampiran"`
	UpdatedAt   string `json:"updatedAt"`
}

// LaporanDispos is LaporanDispo List
type LaporanDispos struct {
	LaporanDispos []LaporanDispo `json:"laporan_dispo"`
}

// GetLaporanDispo is func
func GetLaporanDispo(idDisposisi string) LaporanDispos {
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
	_, err := con.Exec("INSERT INTO laporan_disposisi (idDisposisi, idPenerima, status) VALUES (?,?,?)", laporan.IDDisposisi, laporan.IDPenerima, laporan.Status)
	if err != nil {
		fmt.Println("error1 : ", err.Error())
	}
	_, err = con.Exec("UPDATE laporan_disposisi status = 'Forward' WHERE idDisposisi = ? AND idPenerima = ?", laporan.IDDisposisi, idPemberi)
	if err != nil {
		fmt.Println("error2 : ", err.Error())
	}

	defer con.Close()
}
