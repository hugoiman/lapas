package models

import "lapas/db"

// SubSurvei is class
type SubSurvei struct {
	IDSub     int    `json:"idSub"`
	SubSurvei string `json:"subSurvei" validate:"required,min=3,max=50"`
	Deleted   bool   `json:"deleted,omitempty"`
}

// SubSurveis is list of sub survei
type SubSurveis struct {
	SubSurveis []SubSurvei `json:"subSurvei"`
}

// GetSubSurvei is func
func GetSubSurvei() SubSurveis {
	con := db.Connect()
	query := "SELECT idSub, subSurvei FROM subsurvei WHERE deleted = 0"
	rows, _ := con.Query(query)

	subSurvei := SubSurvei{}
	subSurveis := SubSurveis{}

	for rows.Next() {
		_ = rows.Scan(
			&subSurvei.IDSub, &subSurvei.SubSurvei)

		subSurveis.SubSurveis = append(subSurveis.SubSurveis, subSurvei)
	}

	defer con.Close()
	return subSurveis
}

// CreateSubSurvei is func
func CreateSubSurvei(sub SubSurvei) {
	con := db.Connect()
	_, _ = con.Exec("INSERT INTO subsurvei (subSurvei, deleted) VALUES (?,?)", sub.SubSurvei, sub.Deleted)

	defer con.Close()
}

// DeleteSubSurvei is update status delete
func DeleteSubSurvei(idSub string) bool {
	con := db.Connect()
	query := "UPDATE subsurvei SET deleted = 1 WHERE idSub = ?"
	res, _ := con.Exec(query, idSub)

	count, _ := res.RowsAffected()

	defer con.Close()

	if int(count) == 0 {
		return false
	}

	return true
}
