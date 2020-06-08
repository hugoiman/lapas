package models

import (
	"lapas/db"
	"time"
)

// Surat is class
type Surat struct {
	IDSurat     int    `json:"idSurat"`
	Nomor       string `json:"nomor" validate:"required"`
	Sifat       string `json:"sifat" validate:"required"`
	Status      string `json:"status"`
	Perihal     string `json:"perihal" validate:"required"`
	Asal        string `json:"asal" validate:"required"`
	Tujuan      string `json:"tujuan" validate:"required"`
	Penerima    string `json:"penerima" validate:"required"`
	Lampiran    string `json:"lampiran" validate:"required"`
	InputBy     string `json:"inputBy"`
	UpdatedBy   string `json:"updatedBy"`
	TglSurat    string `json:"tglSurat" validate:"required"`
	TglDiterima string `json:"tglDiterima"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// Surats is surat list
type Surats struct {
	Surats []Surat `json:"surat"`
}

// GetSurat is func
func GetSurat(idSurat string) (Surat, error) {
	con := db.Connect()
	query := "SELECT idSurat, nomor, sifat, status, perihal, asal, tujuan, (SELECT a.nama FROM user a JOIN surat b ON a.idUser = b.penerima WHERE b.idSurat = ?) AS penerima, lampiran, (SELECT a.nama FROM user a JOIN surat b ON a.idUser = b.inputBy WHERE b.idSurat = ?) AS inputBy, (SELECT IFNULL((SELECT a.nama FROM user a JOIN surat b ON a.idUser = b.updatedBy WHERE b.idSurat = ?),'')) AS updatedBy, tglSurat, tglDiterima, createdAt, updatedAt FROM surat WHERE idSurat = ?"

	surat := Surat{}
	var tglSurat, createdAt time.Time
	var updatedAt, tglDiterima interface{}

	err := con.QueryRow(query, idSurat, idSurat, idSurat, idSurat).Scan(
		&surat.IDSurat, &surat.Nomor, &surat.Sifat,
		&surat.Status, &surat.Perihal, &surat.Asal, &surat.Tujuan,
		&surat.Penerima, &surat.Lampiran, &surat.InputBy, &surat.UpdatedBy,
		&tglSurat, &tglDiterima, &createdAt, &updatedAt)

	surat.TglSurat = tglSurat.Format("02 Jan 2006")
	surat.CreatedAt = createdAt.Format("02 Jan 2006")

	if tglDiterima == nil {
		surat.TglDiterima = ""
	} else {
		surat.TglDiterima = tglDiterima.(time.Time).Format("02 Jan 2006")
	}

	if updatedAt == nil {
		surat.UpdatedAt = ""
	} else {
		surat.UpdatedAt = updatedAt.(time.Time).Format("02 Jan 2006")
	}

	defer con.Close()
	return surat, err
}

// GetSurats is func
func GetSurats() Surats {
	con := db.Connect()
	query := "SELECT idSurat, nomor, sifat, status, perihal, asal, tujuan, penerima, lampiran, inputBy, updatedBy, tglSurat, tglDiterima, createdAt, updatedAt FROM surat"
	rows, _ := con.Query(query)

	surat := Surat{}
	surats := Surats{}

	var tglSurat, tglDiterima, createdAt, updatedAt interface{}

	for rows.Next() {
		_ = rows.Scan(
			&surat.IDSurat, &surat.Nomor, &surat.Sifat,
			&surat.Status, &surat.Perihal, &surat.Asal, &surat.Tujuan,
			&surat.Penerima, &surat.Lampiran, &surat.InputBy, &surat.UpdatedBy,
			&tglSurat, &tglDiterima, &createdAt, &updatedAt)

		surat.TglSurat = tglSurat.(time.Time).Format("02 Jan 2006")
		surat.CreatedAt = createdAt.(time.Time).Format("02 Jan 2006")

		if tglDiterima == nil {
			surat.TglDiterima = ""
		} else {
			surat.TglDiterima = tglDiterima.(time.Time).Format("02 Jan 2006")
		}

		if updatedAt == nil {
			surat.UpdatedAt = ""
		} else {
			surat.UpdatedAt = updatedAt.(time.Time).Format("02 Jan 2006")
		}

		surats.Surats = append(surats.Surats, surat)
	}

	defer con.Close()
	return surats
}

// CreateSurat is func
func CreateSurat(surat Surat) error {
	con := db.Connect()

	var err error
	if surat.TglDiterima == "" {
		_, err = con.Exec("INSERT INTO surat (nomor, sifat, status, perihal, asal, tujuan, penerima, lampiran, inputBy, tglSurat, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)", surat.Nomor, surat.Sifat, surat.Status, surat.Perihal, surat.Asal, surat.Tujuan, surat.Penerima, surat.Lampiran, surat.InputBy, surat.TglSurat, surat.CreatedAt)
	} else {
		_, err = con.Exec("INSERT INTO surat (nomor, sifat, status, perihal, asal, tujuan, penerima, lampiran, inputBy, tglSurat, tglDiterima, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", surat.Nomor, surat.Sifat, surat.Status, surat.Perihal, surat.Asal, surat.Tujuan, surat.Penerima, surat.Lampiran, surat.InputBy, surat.TglSurat, surat.TglDiterima, surat.CreatedAt)

	}

	defer con.Close()

	return err
}

// UpdateSurat is func
func UpdateSurat(idSurat string, surat Surat) error {
	con := db.Connect()
	var err error

	if surat.TglDiterima == "" {
		query := "UPDATE surat SET nomor = ?, sifat = ?, status = ?, perihal = ?, asal = ?, tujuan = ?, penerima = ?, lampiran = ?, updatedBy = ?, tglSurat = ?, updatedAt = ? WHERE idSurat = ? AND status != 'Deleted'"
		_, err = con.Exec(query, surat.Nomor, surat.Sifat, surat.Status, surat.Perihal, surat.Asal, surat.Tujuan, surat.Penerima, surat.Lampiran, surat.UpdatedBy, surat.TglSurat, surat.UpdatedAt, idSurat)

	} else {
		query := "UPDATE surat SET nomor = ?, sifat = ?, status = ?, perihal = ?, asal = ?, tujuan = ?, penerima = ?, lampiran = ?, updatedBy = ?, tglSurat = ?, tglDiterima = ?, updatedAt = ? WHERE idSurat = ? AND status != 'Deleted'"
		_, err = con.Exec(query, surat.Nomor, surat.Sifat, surat.Status, surat.Perihal, surat.Asal, surat.Tujuan, surat.Penerima, surat.Lampiran, surat.UpdatedBy, surat.TglSurat, surat.TglDiterima, surat.UpdatedAt, idSurat)
	}

	defer con.Close()

	return err
}

// DeleteSurat is func
func DeleteSurat(idSurat, deletedBy, updated string) {
	con := db.Connect()
	query := "UPDATE surat SET status = 'Deleted', updatedBy = ?, updatedAt = ? WHERE idSurat = ?"
	_, _ = con.Exec(query, deletedBy, updated, idSurat)

	defer con.Close()
}

// BeriStatus is func
func BeriStatus(idSurat string) {
	con := db.Connect()
	query := "UPDATE surat SET status = 'Filling' WHERE idSurat = ?"
	_, _ = con.Exec(query, idSurat)
	defer con.Close()
}
