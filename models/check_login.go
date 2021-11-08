package models

import (
	"database/sql"
	"errors"
)

//登录错误定义
var (
	// ErrorLoginUserNotFound 登录的用户不存在
	ErrorLoginUserNotFound = errors.New("login user not found")
	// ErrorLoginInvalidPassword 密码错误
	ErrorLoginInvalidPassword = errors.New("invalid password")
	// ErrorLoginAccountLocked 账号停权
	ErrorLoginAccountLocked = errors.New("account locked")
	//ErrorLoginAccountOnline 有角色在线
	ErrorLoginAccountOnline = errors.New("account role is online")
)

//CheckLogin 验证登录
func CheckLogin(db *sql.DB, username, password string) error {
	account, err := GetAccountByUsername(db, username)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrorLoginUserNotFound
	}
	if account.Password != password {
		return ErrorLoginInvalidPassword
	}
	if account.Qq.Valid {
		if account.Qq.String == "1" {
			return ErrorLoginAccountLocked
		}
	}
	//todo online
	return nil
}
