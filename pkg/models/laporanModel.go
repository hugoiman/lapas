package models

import (
	"lapas/db"
	"time"
)

// Laporan is class
type Laporan struct {
	IDLaporan    int    `json:"idLaporan"`
	Pengirim     string `json:"pengirim" validate:"required"`
	Pangkat      string `json:"pangkat"`
	Divisi       string `json:"divisi"`
	Subjek       string `json:"subjek" validate:"required"`
	Kategori     string `json:"kategori" validate:"required"`
	Pesan        string `json:"pesan" validate:"required"`
	Lampiran     string `json:"lampiran"`
	Status       string `json:"status"`
	Penanggap    string `json:"penanggap"`
	Tanggapan    string `json:"tanggapan"`
	TglKirim     string `json:"tglKirim"`
	TglTanggapan string `json:"tglTanggapan"`
}

// Laporans is laporan list
type Laporans struct {
	Laporans []Laporan `json:"laporan"`
}

// GetLaporan is func
func GetLaporan(idLaporan string) (Laporan, error) {
	con := db.Connect()
	query := "SELECT a.idLaporan, b.nama, b.pangkat, b.divisi, a.subjek, a.kategori, a.pesan, a.lampiran, a.status, a.penanggap, a.tanggapan, a.tglKirim, a.tglTanggapan FROM laporan a JOIN user b ON a.pengirim = b.idUser WHERE a.idLaporan = ?"

	laporan := Laporan{}
	var tglKirim time.Time

	err := con.QueryRow(query, idLaporan).Scan(
		&laporan.IDLaporan, &laporan.Pengirim, &laporan.Pangkat, &laporan.Divisi, &laporan.Subjek, &laporan.Kategori, &laporan.Pesan,
		&laporan.Lampiran, &laporan.Status, &laporan.Penanggap, &laporan.Tanggapan, &tglKirim, &laporan.TglTanggapan)

	laporan.TglKirim = tglKirim.Format("02 Jan 2006")

	defer con.Close()
	return laporan, err
}

// GetLaporans is func
func GetLaporans() Laporans {
	con := db.Connect()
	query := "SELECT a.idLaporan, b.nama, a.subjek, a.kategori, a.pesan, a.lampiran, a.status, a.tglKirim FROM laporan a JOIN user b ON a.pengirim = b.idUser "
	rows, _ := con.Query(query)

	laporan := Laporan{}
	laporans := Laporans{}
	var tglKirim time.Time

	for rows.Next() {
		_ = rows.Scan(
			&laporan.IDLaporan, &laporan.Pengirim, &laporan.Subjek, &laporan.Kategori, &laporan.Pesan, &laporan.Lampiran, &laporan.Status, &tglKirim)

		laporan.TglKirim = tglKirim.Format("02 Jan 2006")
		laporans.Laporans = append(laporans.Laporans, laporan)
	}

	defer con.Close()
	return laporans
}

// GetMyLaporan is func
func GetMyLaporan(idUser string) Laporans {
	con := db.Connect()
	query := "SELECT idLaporan, subjek, kategori, pesan, lampiran, status, tglKirim FROM laporan WHERE pengirim = ?"
	rows, _ := con.Query(query, idUser)

	laporan := Laporan{}
	laporans := Laporans{}
	var tglKirim time.Time

	for rows.Next() {
		_ = rows.Scan(
			&laporan.IDLaporan, &laporan.Subjek, &laporan.Kategori, &laporan.Pesan, &laporan.Lampiran, &laporan.Status, &tglKirim)

		laporan.TglKirim = tglKirim.Format("02 Jan 2006")
		laporans.Laporans = append(laporans.Laporans, laporan)
	}

	defer con.Close()
	return laporans
}

// CreateLaporan is func
func CreateLaporan(laporan Laporan) error {
	con := db.Connect()

	_, err := con.Exec("INSERT INTO laporan (pengirim, subjek, kategori, pesan, lampiran, status, tglKirim) VALUES (?,?,?,?,?,?,?)", laporan.Pengirim, laporan.Subjek, laporan.Kategori, laporan.Pesan, laporan.Lampiran, laporan.Status, laporan.TglKirim)

	defer con.Close()

	return err
}

// CreateTanggapan is func
func CreateTanggapan(idLaporan string, tanggapan Laporan) int {
	con := db.Connect()
	query := "UPDATE laporan SET tanggapan = ?, penanggap = ?, tglTanggapan = ?, status = ? WHERE idLaporan = ?"
	res, _ := con.Exec(query, tanggapan.Tanggapan, tanggapan.Penanggap, tanggapan.TglTanggapan, tanggapan.Status, idLaporan)

	count, _ := res.RowsAffected()

	defer con.Close()

	return int(count)
}
