package models

import (
	"database/sql"
	"errors"
	"github.com/liuguangw/billing_go/common"
)

// 登录错误定义
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

// CheckLogin 验证登录
func CheckLogin(db *sql.DB, onlineUsers map[string]*common.ClientInfo, username, password string) error {
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
	if account.IDCard.Valid {
		if account.IDCard.String == "1" {
			return ErrorLoginAccountLocked
		}
	}
	//判断用户是否在线
	if _, userOnline := onlineUsers[username]; userOnline {
		return ErrorLoginAccountOnline
	}
	return nil
}
