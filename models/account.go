package models

import (
	"database/sql"
)

//Account 用户信息结构
type Account struct {
	Id       int32
	Name     string
	Password string
	Question sql.NullString
	Answer   sql.NullString
	Email    sql.NullString
	Qq       sql.NullString
	Point    int
}

// UpdateOnlineStatus 更新在线状态
func UpdateOnlineStatus(db *sql.DB, username string, isOnline bool) error {
	stmt, err := db.Prepare("UPDATE account SET is_online=? WHERE name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var onlineStatus byte
	if isOnline {
		onlineStatus = 1
	}
	_, err = stmt.Exec(onlineStatus, username)
	return err
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
