package models

import "database/sql"

// GetAccountByUsername 通过用户名获取用户信息
func GetAccountByUsername(db *sql.DB, username string) (*Account, error) {
	var account Account
	row := db.QueryRow("SELECT id,name,password"+
		",question,answer,email,id_card,point"+
		" FROM account WHERE name=?", username)
	if err := row.Err(); err != nil {
		return nil, err
	}
	if err := row.Scan(&account.ID, &account.Name, &account.Password,
		&account.Question, &account.Answer, &account.Email, &account.IDCard, &account.Point); err != nil {
		if err == sql.ErrNoRows {
			// 查询不到此用户名的记录
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}
