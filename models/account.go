package models

import (
	"database/sql"
	"errors"
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
	IsOnline byte
	IsLock   byte
}

//通过用户名获取用户信息
func GetAccountByUsername(db *sql.DB, username string) (*Account, error) {
	var account Account
	rows, err := db.Query("SELECT id,name,password"+
		",question,answer,email,qq,point"+
		",is_online,is_lock"+
		" FROM account WHERE name=?", username)
	if err != nil {
		return nil, errors.New("db error: " + err.Error())
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&account.Id, &account.Name, &account.Password,
			&account.Question, &account.Answer, &account.Email, &account.Qq, &account.Point,
			&account.IsOnline, &account.IsLock)
		if err != nil {
			return nil, errors.New("db rows.Scan error: " + err.Error())
		}
	} else {
		// 查询不到此用户名的记录
		return nil, nil
	}
	return &account, nil
}

//获取登录结果
func GetLoginResult(db *sql.DB, username, password string) (byte, error) {
	account, err := GetAccountByUsername(db, username)
	if err != nil {
		//数据库异常
		return 6, err
	} else if account == nil {
		// 用户不存在
		return 9, errors.New("user " + username + " does not exists(go to register)")
	}
	if account.Password != password {
		// 密码错误
		return 3, errors.New("user " + username + " password error")
	}
	if account.IsLock != 0 {
		//停权
		return 7, errors.New("user " + username + " account locked")
	}
	if account.IsOnline != 0 {
		//有角色在线
		return 4, errors.New("user " + username + " is online")
	}
	return 1, nil
}

//注册
func GetRegisterResult(db *sql.DB, username string, password string, superPassword string, email string) (byte, error) {
	//成功、失败状态码
	var (
		regSuccessCode byte = 1
		regErrCode     byte = 4
	)
	account, err := GetAccountByUsername(db, username)
	if err != nil {
		//数据库异常
		return regErrCode, err
	} else if account != nil {
		// 要注册的用户已存在
		return regErrCode, errors.New("user " + username + " already exists")
	}
	// 不允许默认的邮箱
	if email == "1@1.com" {
		return regErrCode, errors.New("email " + email + " is not allowed")
	}
	stmt, err := db.Prepare("INSERT INTO account (name, password, question, email) VALUES (?, ?, ?, ?)")
	if err != nil {
		return regErrCode, errors.New("db.Prepare error: " + err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password, superPassword, email)
	if err != nil {
		return regErrCode, errors.New("db stmt.Exec error: " + err.Error())
	}
	return regSuccessCode, nil
}

// 更新在线状态
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

// 点数兑换
func ConvertUserPoint(db *sql.DB, username string, realPoint int) error {
	stmt, err := db.Prepare("UPDATE account SET point=point-? WHERE name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(realPoint, username)
	return err
}
