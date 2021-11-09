package models

import (
	"database/sql"
)

//Account 用户信息结构
type Account struct {
	ID       int32
	Name     string
	Password string
	Question sql.NullString
	Answer   sql.NullString
	Email    sql.NullString
	IDCard   sql.NullString
	Point    int
}

// ConvertUserPoint 点数兑换
func ConvertUserPoint(db *sql.DB, username string, realPoint int) error {
	stmt, err := db.Prepare("UPDATE account SET point=point-? WHERE name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(realPoint, username)
	return err
}
