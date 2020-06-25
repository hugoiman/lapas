package models

import (
	"lapas/db"
	"strconv"
	"time"
)

// Disposisi is class
type Disposisi struct {
	IDDisposisi   int            `json:"idDisposisi"`
	IDSurat       int            `json:"idSurat" validate:"required"`
	Instruksi     string         `json:"instruksi" validate:"required"`
	Status        string         `json:"status"`
	IDPemberi     int            `json:"idPemberi"`
	Pemberi       string         `json:"pemberi"`
	CreatedAt     string         `json:"createdAt"`
	UpdatedAt     string         `json:"updatedAt"`
	LaporanDispos []LaporanDispo `json:"laporan_dispo" validate:"required"`
}

// Disposisis is disposisi list
type Disposisis struct {
	Disposisis []Disposisi `json:"disposisi"`
}

// GetDisposisi is func
func GetDisposisi(idDisposisi string) (Disposisi, error) {
	con := db.Connect()
	query := "SELECT idDisposisi, idSurat, instruksi, status, idPemberi, " +
		"(SELECT nama FROM user WHERE idUser = idPemberi) AS pemberi, createdAt, updatedAt FROM disposisi WHERE idDisposisi = ?"

	disposisi := Disposisi{}
	var createdAt time.Time
	var updatedAt interface{}

	err := con.QueryRow(query, idDisposisi).Scan(
		&disposisi.IDDisposisi, &disposisi.IDSurat, &disposisi.Instruksi,
		&disposisi.Status, &disposisi.IDPemberi, &disposisi.Pemberi, &createdAt, &updatedAt)

	disposisi.CreatedAt = createdAt.Format("02 Jan 2006")
	if updatedAt == nil {
		disposisi.UpdatedAt = ""
	} else {
		disposisi.UpdatedAt = updatedAt.(time.Time).Format("02 Jan 2006")
	}

	laporan := GetLaporanDispos(idDisposisi)
	disposisi.LaporanDispos = laporan.LaporanDispos

	defer con.Close()
	return disposisi, err
}

// GetDispoSurat is func
func GetDispoSurat(idSurat string) Disposisis {
	con := db.Connect()
	query := "SELECT idDisposisi, idSurat, instruksi, status, idPemberi, " +
		"(SELECT nama FROM user WHERE idUser = idPemberi) AS pemberi, createdAt, updatedAt FROM disposisi WHERE idSurat = ?"
	rows, _ := con.Query(query, idSurat)

	disposisi := Disposisi{}
	disposisis := Disposisis{}

	var createdAt time.Time
	var updatedAt interface{}

	for rows.Next() {
		_ = rows.Scan(
			&disposisi.IDDisposisi, &disposisi.IDSurat, &disposisi.Instruksi,
			&disposisi.Status, &disposisi.IDPemberi, &disposisi.Pemberi, &createdAt, &updatedAt)

		disposisi.CreatedAt = createdAt.Format("02 Jan 2006")
		if updatedAt == nil {
			disposisi.UpdatedAt = ""
		} else {
			disposisi.UpdatedAt = updatedAt.(time.Time).Format("02 Jan 2006")
		}

		disposisis.Disposisis = append(disposisis.Disposisis, disposisi)
	}

	// Get laporan dispo
	for k, v := range disposisis.Disposisis {
		laporan := GetLaporanDispos(strconv.Itoa(v.IDDisposisi))
		disposisis.Disposisis[k].LaporanDispos = laporan.LaporanDispos
	}

	defer con.Close()
	return disposisis
}

// GetDisposisis is func
func GetDisposisis() Surats {
	con := db.Connect()
	query := "SELECT idSurat, nomor, sifat, status, perihal, asal, tujuan, idPenerima," +
		"(SELECT nama FROM user WHERE idUser = idPenerima) AS penerima, lampiran, inputById," +
		"(SELECT nama FROM user WHERE idUser = inputById) AS inputBy, updatedById," +
		"(SELECT IFNULL((SELECT nama FROM user WHERE idUser = updatedById),'')) AS updatedBy," +
		"tglSurat, tglDiterima, createdAt, updatedAt FROM surat WHERE status = 'Dispo'"
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

// GetMyDisposisis is func
func GetMyDisposisis(idUser string) Surats {
	con := db.Connect()
	query := "SELECT a.idSurat, a.nomor, a.status, a.perihal, a.asal, a.tujuan, a.tglSurat FROM surat a JOIN disposisi b ON a.idSurat = b.idSurat JOIN laporan_disposisi c ON b.idDisposisi = c.idDisposisi WHERE c.penerima = ?"
	rows, _ := con.Query(query, idUser)

	surat := Surat{}
	surats := Surats{}

	var tglSurat time.Time

	for rows.Next() {
		_ = rows.Scan(&surat.IDSurat, &surat.Nomor, &surat.Status, &surat.Perihal, &surat.Asal, &surat.Tujuan, &tglSurat)

		surat.TglSurat = tglSurat.Format("02 Jan 2006")
		surats.Surats = append(surats.Surats, surat)
		surat.Disposisis = make([]Disposisi, 0)
	}

	for k, v := range surats.Surats {
		query := "SELECT idDisposisi, instruksi, status, pemberi FROM disposisi WHERE idSurat = ?"
		disposisi := Disposisi{}
		_ = con.QueryRow(query, v.IDSurat).Scan(&disposisi.IDDisposisi, &disposisi.Instruksi, &disposisi.Status, &disposisi.Pemberi)
		surats.Surats[k].Disposisis = append(surats.Surats[k].Disposisis, disposisi)
	}

	defer con.Close()
	return surats
}

// CreateDisposisi is func
func CreateDisposisi(disposisi Disposisi) int {
	con := db.Connect()
	exec, _ := con.Exec("INSERT INTO disposisi (idSurat, instruksi, status, pemberi, createdAt) VALUES (?,?,?,?,?)",
		disposisi.IDSurat, disposisi.Instruksi, disposisi.Status, disposisi.Pemberi, disposisi.CreatedAt)

	idInt64, _ := exec.LastInsertId()
	idDisposisi := int(idInt64)
	defer con.Close()

	return idDisposisi
}

// BeriStatusDisposisi is func
func BeriStatusDisposisi(idDisposisi, updatedByID, updatedAt, status string) {
	con := db.Connect()
	var query string

	if status == "Solved" {
		query = "UPDATE disposisi a JOIN surat b ON a.idSurat = b.idSurat " +
			"SET a.status = 'Solved', a.updatedAt = ?, b.status = 'Solved', b.updatedById = ?, b.updatedAt = ? WHERE a.idDisposisi = ?"
		_, _ = con.Exec(query, idDisposisi)
	} else if status == "Waiting" {
		// hapus laporan yang id disposisi = ?
		query = "DELETE FROM laporan_disposisi WHERE idDisposisi = ?"
		_, _ = con.Exec(query, idDisposisi)
	}

	defer con.Close()
}
