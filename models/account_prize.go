package models

import (
	"database/sql"
	"strings"
)

// AccountPrize 对应奖励表
//
// 创建表的SQL语句
/*
CREATE TABLE account_prize (
  id bigint(20) NOT NULL AUTO_INCREMENT,
  account varchar(50) NOT NULL COMMENT '账号',
  world int(11) NOT NULL DEFAULT '0' COMMENT '世界ID',
  charguid int(10) unsigned NOT NULL DEFAULT '0' COMMENT '玩家GUID',
  itemid int(10) unsigned NOT NULL DEFAULT '0' COMMENT '物品ID',
  itemnum int(11) NOT NULL COMMENT '物品数量',
  isget smallint(6) NOT NULL COMMENT '是否领取了',
  validtime int(11) NOT NULL COMMENT '有效期，时间格式为unix时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY id (id) USING BTREE
)
*/
type AccountPrize struct {
	ID        int64
	Account   string
	World     int
	Charguid  int
	ItemID    int
	ItemNum   int
	IsGet     int
	ValidTime int
}

// CheckAccountPrizeState 查询奖励状态
func CheckAccountPrizeState(db *sql.DB, username string, world, charguid int) (byte, error) {
	var state byte
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM account_prize WHERE account=? AND world=? AND charguid=? AND isget=0) as m", username, world, charguid)
	if err := row.Err(); err != nil {
		return 0, err
	}
	if err := row.Scan(&state); err != nil {
		return 0, err
	}
	return state, nil
}

// FindFirstAccountPrize 查询当前用户第一条奖励记录
func FindFirstAccountPrize(db *sql.DB, username string) (*AccountPrize, error) {
	var accountPrize AccountPrize
	sqlStr := "SELECT id,account,world,charguid,itemid,itemnum,isget,validtime" +
		" FROM account_prize WHERE account=? AND isget=0" +
		" ORDER BY id LIMIT 1"
	row := db.QueryRow(sqlStr, username)
	if err := row.Err(); err != nil {
		return nil, err
	}
	if err := row.Scan(&accountPrize.ID, &accountPrize.Account, &accountPrize.World, &accountPrize.Charguid,
		&accountPrize.ItemID, &accountPrize.ItemNum, &accountPrize.IsGet, &accountPrize.ValidTime); err != nil {
		if err == sql.ErrNoRows {
			// 查询不到此记录
			return nil, nil
		}
		return nil, err
	}
	return &accountPrize, nil
}

// FindAccountPrizeList 查询符合条件的奖励项列表
func FindAccountPrizeList(db *sql.DB, username string, world, charguid, limitCount int) ([]AccountPrize, error) {
	var prizeList []AccountPrize
	sqlStr := "SELECT id,account,world,charguid,itemid,itemnum,isget,validtime" +
		" FROM account_prize WHERE account=? AND world=? AND charguid=? AND isget=0" +
		" ORDER BY id LIMIT ?"
	rows, err := db.Query(sqlStr, username, world, charguid, limitCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var accountPrize AccountPrize
		if err := rows.Scan(&accountPrize.ID, &accountPrize.Account, &accountPrize.World, &accountPrize.Charguid,
			&accountPrize.ItemID, &accountPrize.ItemNum, &accountPrize.IsGet, &accountPrize.ValidTime); err != nil {
			return nil, err
		}
		prizeList = append(prizeList, accountPrize)
	}
	return prizeList, nil
}

// MarkGetAccountPrize 标记为已领取奖励
func MarkGetAccountPrize(db *sql.DB, itemIdList []int64) error {
	sqlStr := "UPDATE account_prize SET isget=1 WHERE id IN"
	params := make([]string, len(itemIdList))
	args := make([]any, len(itemIdList))
	for i, id := range itemIdList {
		params[i] = "?"
		args[i] = id
	}
	placeholders := strings.Join(params, ",")
	sqlStr += " (" + placeholders + ")"
	//exec
	_, err := db.Exec(sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}
