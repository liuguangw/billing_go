package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liuguangw/billing_go/config"
	"time"
)

// 获取数据库连接 返回连接对象、版本、error
func GetConnection(sConfig *config.ServerConfig) (db *sql.DB, dbVersion string, err error) {
	//user:password@tcp(localhost:3306)/dbname?charset=utf8....
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sConfig.DbUser, sConfig.DbPassword,
		sConfig.DbHost, sConfig.DbPort, sConfig.DbName)
	extraParams := "?charset=utf8&timeout=1s&readTimeout=1s&writeTimeout=1s"
	if sConfig.AllowOldPassword {
		extraParams += "&allowOldPasswords=true"
	}
	db, err = sql.Open("mysql", connString+extraParams)
	if err != nil {
		return
	}
	// 最大100个连接，最多闲置10个
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	// 连接池连接存活时长
	maxLifeTime, err := time.ParseDuration("30m")
	if err != nil {
		return
	}
	db.SetConnMaxLifetime(maxLifeTime)
	// 判断连接状态
	err = db.Ping()
	if err != nil {
		return
	}
	//获取版本信息
	rows, err := db.Query("SELECT VERSION() as v")
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&dbVersion)
		if err != nil {
			return
		}
	} else {
		return
	}
	// 额外字段的存在性初始化
	extraFields := map[string]bool{
		"is_online": false,
		"is_lock":   false,
	}
	// 获取account表的所以字段信息
	rows, err = db.Query("SHOW COLUMNS FROM account")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			fieldName string
			tType     string
			tNull     string
			tKey      string
			tDefault  sql.NullString
			tExtra    string
		)
		tErr := rows.Scan(&fieldName, &tType, &tNull, &tKey, &tDefault, &tExtra)
		if tErr != nil {
			err = tErr
			return
		}
		for tFieldName, tFieldExists := range extraFields {
			if !tFieldExists && tFieldName == fieldName {
				// 标记字段为已存在
				extraFields[tFieldName] = true
				break
			}
		}
		//fmt.Printf("fieldName: %s\n",fieldName)
	}
	// 需要添加的额外字段名
	var needAddFields []string
	for tFieldName, tFieldExists := range extraFields {
		if !tFieldExists {
			needAddFields = append(needAddFields, tFieldName)
		}
	}
	if len(needAddFields) > 0 {
		// 添加额外字段
		for _, fName := range needAddFields {
			stmt, tErr := db.Prepare("ALTER TABLE account ADD COLUMN " + fName + " smallint(1) UNSIGNED NOT NULL DEFAULT 0")
			if tErr != nil {
				err = tErr
				return
			}
			_, tErr = stmt.Exec()
			_ = stmt.Close()
			if tErr != nil {
				err = tErr
				return
			}
		}
	}
	return
}
