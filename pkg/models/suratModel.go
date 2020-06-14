package models

import (
	"lapas/db"
	"strconv"
	"time"
)

// Surat is class
type Surat struct {
	IDSurat     int         `json:"idSurat"`
	Nomor       string      `json:"nomor" validate:"required,min=4,max=50"`
	Sifat       string      `json:"sifat" validate:"required,eq=Biasa|eq=Segera|eq=Penting"`
	Status      string      `json:"status"`
	Perihal     string      `json:"perihal" validate:"required,max=50"`
	Asal        string      `json:"asal" validate:"required,max=50"`
	Tujuan      string      `json:"tujuan" validate:"required,max=50"`
	IDPenerima  int         `json:"idPenerima" validate:"required"`
	Penerima    string      `json:"penerima"`
	Lampiran    string      `json:"lampiran" validate:"required,endswith=.docx|endswith=.pdf"`
	InputByID   int         `json:"inputById"`
	InputBy     string      `json:"inputBy"`
	UpdatedByID int         `json:"updatedById"`
	UpdatedBy   string      `json:"updatedBy"`
	TglSurat    string      `json:"tglSurat" validate:"required"`
	TglDiterima string      `json:"tglDiterima"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
	Disposisis  []Disposisi `json:"disposisi"`
}

// Surats is surat list
type Surats struct {
	Surats []Surat `json:"surat"`
}

// GetSurat is func
func GetSurat(idSurat string) (Surat, error) {
	con := db.Connect()
	query := "SELECT idSurat, nomor, sifat, status, perihal, asal, tujuan, idPenerima," +
		"(SELECT a.nama FROM user a JOIN surat b ON a.idUser = b.idPenerima WHERE b.idSurat = ?) AS penerima, lampiran, inputById," +
		"(SELECT a.nama FROM user a JOIN surat b ON a.idUser = b.inputById WHERE b.idSurat = ?) AS inputBy, updatedById," +
		"(SELECT IFNULL((SELECT a.nama FROM user a JOIN surat b ON a.idUser = b.updatedById WHERE b.idSurat = ?),'')) AS updatedBy," +
		"tglSurat, tglDiterima, createdAt, updatedAt FROM surat WHERE idSurat = ?"

	surat := Surat{}
	var tglSurat, createdAt time.Time
	var updatedAt, tglDiterima interface{}

	err := con.QueryRow(query, idSurat, idSurat, idSurat, idSurat).Scan(
		&surat.IDSurat, &surat.Nomor, &surat.Sifat,
		&surat.Status, &surat.Perihal, &surat.Asal, &surat.Tujuan,
		&surat.IDPenerima, &surat.Penerima, &surat.Lampiran,
		&surat.InputByID, &surat.InputBy, &surat.UpdatedByID, &surat.UpdatedBy,
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

	dispo := GetDispoSurat(strconv.Itoa(surat.IDSurat))
	surat.Disposisis = dispo.Disposisis

	defer con.Close()
	return surat, err
}

// GetSurats is func
func GetSurats() Surats {
	con := db.Connect()
	query := "SELECT idSurat, nomor, sifat, status, perihal, asal, tujuan, idPenerima," +
		"(SELECT nama FROM user WHERE idUser = idPenerima) AS penerima, lampiran, inputById," +
		"(SELECT nama FROM user WHERE idUser = inputById) AS inputBy, updatedById," +
		"(SELECT IFNULL((SELECT nama FROM user WHERE idUser = updatedById),'')) AS updatedBy," +
		"tglSurat, tglDiterima, createdAt, updatedAt FROM surat"
	rows, _ := con.Query(query)

	surat := Surat{}
	surats := Surats{}

	var tglSurat, tglDiterima, createdAt, updatedAt interface{}

	for rows.Next() {
		_ = rows.Scan(
			&surat.IDSurat, &surat.Nomor, &surat.Sifat,
			&surat.Status, &surat.Perihal, &surat.Asal, &surat.Tujuan,
			&surat.IDPenerima, &surat.Penerima, &surat.Lampiran,
			&surat.InputByID, &surat.InputBy, &surat.UpdatedByID, &surat.UpdatedBy,
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
		_, err = con.Exec("INSERT INTO surat (nomor, sifat, status, perihal, asal, tujuan, idPenerima, lampiran, inputById, tglSurat, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
			surat.Nomor, surat.Sifat, surat.Status, surat.Perihal, surat.Asal, surat.Tujuan, surat.IDPenerima, surat.Lampiran, surat.InputByID, surat.TglSurat, surat.CreatedAt)
	} else {
		_, err = con.Exec("INSERT INTO surat (nomor, sifat, status, perihal, asal, tujuan, idPenerima, lampiran, inputById, tglSurat, tglDiterima, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
			surat.Nomor, surat.Sifat, surat.Status, surat.Perihal, surat.Asal, surat.Tujuan, surat.IDPenerima, surat.Lampiran, surat.InputByID, surat.TglSurat, surat.TglDiterima, surat.CreatedAt)
	}

	defer con.Close()

	return err
}

// UpdateSurat is func
func UpdateSurat(idSurat string, surat Surat) error {
	con := db.Connect()
	var err error

	if surat.TglDiterima == "" {
		query := "UPDATE surat SET nomor = ?, sifat = ?, perihal = ?, asal = ?, tujuan = ?, idPenerima = ?, lampiran = ?, updatedById = ?, tglSurat = ?, updatedAt = ? WHERE idSurat = ?"
		_, err = con.Exec(query, surat.Nomor, surat.Sifat, surat.Perihal, surat.Asal, surat.Tujuan, surat.IDPenerima, surat.Lampiran, surat.UpdatedByID, surat.TglSurat, surat.UpdatedAt, idSurat)

	} else {
		query := "UPDATE surat SET nomor = ?, sifat = ?, perihal = ?, asal = ?, tujuan = ?, idPenerima = ?, lampiran = ?, updatedById = ?, tglSurat = ?, tglDiterima = ?, updatedAt = ? WHERE idSurat = ?"
		_, err = con.Exec(query, surat.Nomor, surat.Sifat, surat.Perihal, surat.Asal, surat.Tujuan, surat.IDPenerima, surat.Lampiran, surat.UpdatedByID, surat.TglSurat, surat.TglDiterima, surat.UpdatedAt, idSurat)
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

// BeriStatusSurat is func
func BeriStatusSurat(idSurat, updatedByID, updatedAt, status string) {
	con := db.Connect()
	var query string

	if status == "Undelete" {
		query = "UPDATE surat SET status = 'Waiting', updatedById = ?, updatedAt = ? WHERE idSurat = ?"
		_, _ = con.Exec(query, updatedByID, updatedAt, idSurat)

	} else if status == "Filling" {
		query = "UPDATE surat SET status = 'Filing', updatedById = ?, updatedAt = ? WHERE idSurat = ?"
		_, _ = con.Exec(query, updatedByID, updatedAt, idSurat)

		query = "DELETE FROM disposisi WHERE idSurat = ?"
		_, _ = con.Exec(query, idSurat)

	} else if status == "Dispo" {
		query = "UPDATE surat SET status = 'Dispo', updatedById = ?, updatedAt = ? WHERE idSurat = ?"
		_, _ = con.Exec(query, updatedByID, updatedAt, idSurat)
	}

	defer con.Close()
}
