package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateMysql(ip string, port string, userName string, password string, dbName string) (*gorm.DB, error) {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		userName, password, ip, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn))
	return db, err
}
