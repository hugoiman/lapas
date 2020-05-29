package models

import "lapas/db"

// SubSurvei is class
type SubSurvei struct {
	IDSub     int    `json:"idSub"`
	SubSurvei string `json:"subSurvei"`
	Deleted   bool   `json:"deleted"`
}

// SubSurveis is list of sub survei
type SubSurveis struct {
	SubSurveis []SubSurvei `json:"subSurvei"`
}

// GetSubSurvei is func
func GetSubSurvei() SubSurveis {
	con := db.Connect()
	query := ""
	rows, _ := con.Query(query)

	subSurvei := SubSurvei{}
	subSurveis := SubSurveis{}

	for rows.Next() {
		_ = rows.Scan(
			&subSurvei.IDSub, &subSurvei.SubSurvei, &subSurvei.Deleted)
	}

	defer con.Close()
	return subSurveis
}
