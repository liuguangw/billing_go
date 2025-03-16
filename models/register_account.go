package models

import (
	"database/sql"
	"errors"
)

// RegisterAccount 注册账号
func RegisterAccount(db *sql.DB, account *Account) error {
	if len(account.Name) > 50 {
		return errors.New("account name length exceeds 50 characters")
	}
	tmpAccount, err := GetAccountByUsername(db, account.Name)
	if err != nil {
		//数据库异常
		return err
	}
	// 要注册的用户已存在
	if tmpAccount != nil {
		return errors.New("user " + account.Name + " already exists")
	}
	// 邮箱不能为空
	if !account.Email.Valid {
		return errors.New("email can't be empty")
	}
	//不允许默认的邮箱
	if account.Email.String == "1@1.com" {
		return errors.New("email " + account.Email.String + " is not allowed")
	}
	//插入
	stmt, err := db.Prepare("INSERT INTO account (name, password, question, email) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(account.Name, account.Password, account.Question, account.Email); err != nil {
		return err
	}
	return nil
}
