package models

import "database/sql"

// CheckIsGm 查询角色是否为gm,怀旧专用
func CheckIsGm(db *sql.DB, charguid int) (bool, error) {
	var accountCfg AccountCfg
	row := db.QueryRow("SELECT charguid, isgm FROM account_cfg WHERE charguid=?", charguid)
	if err := row.Err(); err != nil {
		return false, err
	}
	if err := row.Scan(&accountCfg.Charguid, &accountCfg.Isgm); err != nil {
		if err == sql.ErrNoRows {
			// 查询不到此记录
			return false, nil
		}
		return false, err
	}
	return accountCfg.Isgm > 0, nil
}
