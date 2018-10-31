package main

import (
	"database/sql"
)

//Account 用户信息结构
type Account struct {
	id       int32
	name     string
	password string
	question sql.NullString
	answer   sql.NullString
	email    sql.NullString
	qq       sql.NullString
	point    int32
	isOnline byte
	isLock   byte
}

// 第二个返回值 0表示查询不到此用户名的记录 1表示查询成功 2表示数据库异常
func getAccountByUsername(db *sql.DB, username string) (*Account, byte) {
	var account Account
	rows, err := db.Query("SELECT id,name,password"+
		",question,answer,email,qq,point,is_online,is_lock"+
		" FROM account WHERE name=?", username)
	if err != nil {
		return &account, 2
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&account.id, &account.name, &account.password,
			&account.question, &account.answer, &account.email, &account.qq, &account.point,
			&account.isOnline, &account.isLock)
		if err != nil {
			return &account, 2
		}
	} else {
		// 查询不到此用户名的记录
		return &account, 0
	}
	return &account, 1
}

func getLoginResult(db *sql.DB, username string, password string) byte {
	account, queryOp := getAccountByUsername(db, username)
	if queryOp == 0 {
		// 用户不存在
		return 9
	} else if queryOp == 2 {
		// 数据库异常
		return 6
	}
	if account.password != password {
		// 密码错误
		return 3
	}
	if account.isLock != 0 {
		//停权
		return 7
	}
	if account.isOnline != 0 {
		//有角色在线
		return 4
	}
	return 1
}

func getRegisterResult(db *sql.DB, username string, password string, superPassword string, email string) byte {
	_, queryOp := getAccountByUsername(db, username)
	var regErr byte = 4
	if queryOp == 1 {
		// 用户已存在
		return regErr
	} else if queryOp == 2 {
		// 数据库异常
		return regErr
	}
	stmt, err := db.Prepare("INSERT INTO account (name, password, question, email) VALUES (?, ?, ?, ?)")
	if err != nil {
		return regErr
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password, superPassword, email)
	if err != nil {
		return regErr
	}
	return 1
}

// 更新在线状态
func updateOnlineStatus(db *sql.DB, username string, isOnline bool) error {
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
