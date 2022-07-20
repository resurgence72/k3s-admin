package db

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k3s-admin/config"
)

var (
	GORM *sql.DB
	err  error
)

func Init() {
	// 组装连接配置
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?chartset=utf8&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPwd,
		config.BbHost,
		config.DbPort,
		"test",
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("db failed")
	}

	sqlDB.SetConnMaxIdleTime(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifeTime)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	GORM = sqlDB

}


func Close() {
	GORM.Close()
}