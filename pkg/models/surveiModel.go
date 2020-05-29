package models

import (
	"fmt"
	"lapas/db"
	"time"
)

// Survei is Class
type Survei struct {
	IDSurvei int       `json:"idSurvei"`
	Judul    string    `json:"judul" validate:"required"`
	Periode  time.Time `json:"periode" validate:"required"`
	Actived  bool      `json:"actived,omitempty"`
	Slug     string    `json:"slug"`
	Soal     []Soal    `json:"soal"`
}

// Surveis is Survei List
type Surveis struct {
	Surveis []Survei `json:"survei"`
}

// GetSurvei is function
func GetSurvei(id string) (Survei, error) {
	con := db.Connect()
	querySurvei := "SELECT idSurvei, judul, periode, actived, slug FROM survei WHERE idSurvei = ? OR slug = ?"

	survei := Survei{}
	var soal Soal
	err := con.QueryRow(querySurvei, id, id).Scan(
		&survei.IDSurvei, &survei.Judul, &survei.Periode, &survei.Actived, &survei.Slug)

	if err != nil {
		return survei, err
	}

	//	GetSoal
	querySoal := "SELECT a.idSoal, a.idSub, a.soal, b.subsurvei FROM soal a JOIN subsurvei b ON a.idSub = b.idSub WHERE a.idSurvei = ?"
	rows, _ := con.Query(querySoal, survei.IDSurvei)

	for rows.Next() {
		_ = rows.Scan(&soal.IDSoal, &soal.IDSub, &soal.Soal, &soal.Subsurvei)
		survei.Soal = append(survei.Soal, soal)
	}

	defer con.Close()

	return survei, err
}

// GetSurveiActived is function
func GetSurveiActived(slug string) (Survei, error) {
	con := db.Connect()
	querySurvei := "SELECT idSurvei, judul, periode, actived, slug FROM survei WHERE slug = ? AND actived = 1"

	survei := Survei{}
	var soal Soal
	err := con.QueryRow(querySurvei, slug).Scan(
		&survei.IDSurvei, &survei.Judul, &survei.Periode, &survei.Actived, &survei.Slug)

	if err != nil {
		return survei, err
	}

	//	GetSoal
	querySoal := "SELECT a.idSoal, a.idSub, a.soal, b.subsurvei FROM soal a JOIN subsurvei b ON a.idSub = b.idSub WHERE a.idSurvei = ?"
	rows, _ := con.Query(querySoal, survei.IDSurvei)

	for rows.Next() {
		_ = rows.Scan(&soal.IDSoal, &soal.IDSub, &soal.Soal, &soal.Subsurvei)
		survei.Soal = append(survei.Soal, soal)
	}

	defer con.Close()

	return survei, err
}

// GetSurveis is function
func GetSurveis() Surveis {
	con := db.Connect()
	query := "SELECT idSurvei, judul, periode, actived, slug FROM survei"
	rows, _ := con.Query(query)

	survei := Survei{}
	surveis := Surveis{}

	for rows.Next() {
		err := rows.Scan(
			&survei.IDSurvei, &survei.Judul, &survei.Periode, &survei.Actived, &survei.Slug)
		if err != nil {
			fmt.Println(err.Error())
		}
		surveis.Surveis = append(surveis.Surveis, survei)
	}

	defer con.Close()
	return surveis
}

// CreateSurvei is new survei
func CreateSurvei(survei Survei) (int, error) {
	con := db.Connect()
	exec, err := con.Exec("INSERT INTO survei (judul, periode, actived, slug) VALUES (?,?,?,?)", survei.Judul, survei.Periode, survei.Actived, survei.Slug)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idSurvei := int(idInt64)
	defer con.Close()

	return idSurvei, err
}

// DeleteSurvei is delete survei
func DeleteSurvei(idSurvei string) int {
	con := db.Connect()
	query := "DELETE FROM survei WHERE idSurvei = ?"
	res, _ := con.Exec(query, idSurvei)

	count, _ := res.RowsAffected()

	defer con.Close()

	return int(count)
}

// UpdateSurvei is edit survei
func UpdateSurvei(idSurvei string, survei Survei) error {
	con := db.Connect()
	query := "UPDATE survei SET judul = ?, periode = ?, actived = ?, slug = ? WHERE idSurvei = ?"
	_, err := con.Exec(query, survei.Judul, survei.Periode, survei.Actived, survei.Slug, idSurvei)

	defer con.Close()

	return err
}

// ChangeStatus is func
func ChangeStatus(idSurvei string, actived bool) {
	con := db.Connect()
	query := "UPDATE survei SET actived = ? WHERE idSurvei = ?"
	_, _ = con.Exec(query, actived, idSurvei)

	defer con.Close()
}

// GetDataResponden is func
func GetDataResponden(idSurvei string) Users {
	con := db.Connect()
	query := "SELECT a.idUser, a.nama, a.pangkat, a.direktorat, a.tglLahir FROM user a JOIN jawaban b ON a.idUser = b.idUser JOIN soal c ON b.idSoal = c.idSoal WHERE c.idSurvei = ? GROUP BY a.idUser"
	rows, err := con.Query(query, idSurvei)

	if err != nil {
		fmt.Println("error:", err.Error())
	}

	user := User{}
	users := Users{}

	for rows.Next() {
		err = rows.Scan(
			&user.IDUser, &user.Nama, &user.Pangkat, &user.Direktorat, &user.TglLahir)
		users.Users = append(users.Users, user)

		if err != nil {
			fmt.Println("error2:", err.Error())
		}
	}

	defer con.Close()
	return users
}
