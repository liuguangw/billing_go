package models

import "database/sql"

// ConvertUserPoint 点数兑换
func ConvertUserPoint(db *sql.DB, username string, realPoint int) error {
	stmt, err := db.Prepare("UPDATE account SET point=point-? WHERE name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(realPoint, username); err != nil {
		return err
	}
	return nil
}
