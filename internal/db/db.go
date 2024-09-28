package db

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"health-server/config"
)

var DB *gorm.DB   // GORM 的 DB 实例
var sqlDB *sql.DB // 原生 SQL DB 实例

// InitDB 初始化数据库连接
func InitDB(dbconfig config.DbConfig) error {
	dsn := dbconfig.User + ":" + dbconfig.Password + "@tcp(" + dbconfig.Addr + ")/" + dbconfig.Dbname + "?charset=utf8&parseTime=true"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err = DB.DB() // 使用 DB 获取 *sql.DB 实例
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(dbconfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbconfig.MaxIdleConns)

	return nil
}
