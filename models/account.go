package models

import "database/sql"

//Account 用户信息结构
type Account struct {
	Id       int32
	Name     string
	Password string
	Question sql.NullString
	Answer   sql.NullString
	Email    sql.NullString
	Qq       sql.NullString
	Point    int32
	IsOnline byte
	IsLock   byte
}

const (
	UserNotFound byte = 0
	UserFound    byte = 1
	DbError      byte = 2
)

// 第二个返回值表示查找状态
func GetAccountByUsername(db *sql.DB, username string) (*Account, byte) {
	var account Account
	rows, err := db.Query("SELECT id,name,password"+
		",question,answer,email,qq,point"+
		",is_online,is_lock"+
		" FROM account WHERE name=?", username)
	if err != nil {
		return &account, DbError
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&account.Id, &account.Name, &account.Password,
			&account.Question, &account.Answer, &account.Email, &account.Qq, &account.Point,
			&account.IsOnline, &account.IsLock)
		if err != nil {
			return &account, DbError
		}
	} else {
		// 查询不到此用户名的记录
		return &account, UserNotFound
	}
	return &account, UserFound
}

//登录
func GetLoginResult(db *sql.DB, username string, password string) byte {
	account, queryOp := GetAccountByUsername(db, username)
	if queryOp == UserNotFound {
		// 用户不存在
		return 9
	} else if queryOp == DbError {
		// 数据库异常
		return 6
	}
	if account.Password != password {
		// 密码错误
		return 3
	}
	if account.IsLock != 0 {
		//停权
		return 7
	}
	if account.IsOnline != 0 {
		//有角色在线
		return 4
	}
	return 1
}

//注册
func GetRegisterResult(db *sql.DB, username string, password string, superPassword string, email string) byte {
	_, queryOp := GetAccountByUsername(db, username)
	var regErr byte = 4
	if queryOp == UserFound {
		// 用户已存在
		return regErr
	} else if queryOp == DbError {
		// 数据库异常
		return regErr
	}
	// 不允许默认的邮箱
	if email == "1@1.com" {
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

