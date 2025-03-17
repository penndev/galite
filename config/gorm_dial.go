package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// 数据库连接实例，
// 将多个数据库连接进行抽象，减少不同数据库配置依赖
// 不论什么数据库最终提供给gram的都为 gorm.Dialector
func GormDial() gorm.Dialector {
	var GormDial gorm.Dialector
	switch {
	case strings.HasPrefix(os.Getenv("DB_URL"), "mariadb://"):
		GormDial = mysql.Open(strings.TrimPrefix(os.Getenv("DB_URL"), "mariadb://"))
	case strings.HasPrefix(os.Getenv("DB_URL"), "mysql://"):
		GormDial = mysql.Open(strings.TrimPrefix(os.Getenv("DB_URL"), "mysql://"))
	case strings.HasPrefix(os.Getenv("DB_URL"), "sqlite://"):
		GormDial = sqlite.Open(strings.TrimPrefix(os.Getenv("DB_URL"), "sqlite://"))
	case strings.HasPrefix(os.Getenv("DB_URL"), "postgres://"):
		GormDial = postgres.Open(os.Getenv("DB_URL"))
	case strings.HasPrefix(os.Getenv("DB_URL"), "sqlserver://"):
		GormDial = sqlserver.Open(os.Getenv("DB_URL"))
	default:
		log.Panic(errors.New("env DB_URL err"))
	}
	return GormDial

}
