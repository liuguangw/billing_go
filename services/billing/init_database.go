package billing

import (
	"context"
	"database/sql"
	"fmt"
	//MySQL 驱动
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// initDatabase 初始化数据库连接
func (s *Server) initDatabase() error {
	//user:password@tcp(localhost:3306)/dbname?charset=utf8....
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", s.config.DbUser, s.config.DbPassword,
		s.config.DbHost, s.config.DbPort, s.config.DbName)
	extraParams := "?charset=utf8"
	if s.config.AllowOldPassword {
		extraParams += "&allowOldPasswords=true"
	}
	db, err := sql.Open("mysql", connString+extraParams)
	if err != nil {
		return err
	}
	//连接最长存活时间
	db.SetConnMaxLifetime(time.Minute * 4)
	// 最大100个连接，最多闲置10个
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	// 判断连接状态
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	//获取版本信息
	var dbVersion string
	row := db.QueryRow("SELECT VERSION() as v")
	if err := row.Err(); err != nil {
		return err
	}
	if err := row.Scan(&dbVersion); err != nil {
		return err
	}
	s.logger.Info("MySQL version: " + dbVersion)
	s.database = db
	return nil
}
