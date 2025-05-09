package models

import (
	"database/sql"
)

// Account 用户信息结构
type Account struct {
	ID       int32          //账号id
	Name     string         //用户名
	Password string         //已加密的密码
	Question sql.NullString //密保问题(现在存放的是加密的超级密码)
	Answer   sql.NullString //密保答案
	Email    sql.NullString //注册时填写的邮箱
	IDCard   sql.NullString //用于判定账号是否锁定,如果是字符串1就说明账号已锁定
	Point    int            //点数
}
