package model

import (
	"errors"
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
func Dial(dsn string) (gorm.Dialector, error) {
	var GormDial gorm.Dialector
	switch {
	case strings.HasPrefix(dsn, "mariadb://"):
		GormDial = mysql.Open(strings.TrimPrefix(dsn, "mariadb://"))
	case strings.HasPrefix(dsn, "mysql://"):
		GormDial = mysql.Open(strings.TrimPrefix(dsn, "mysql://"))
	case strings.HasPrefix(dsn, "sqlite://"):
		GormDial = sqlite.Open(strings.TrimPrefix(dsn, "sqlite://"))
	case strings.HasPrefix(dsn, "postgres://"):
		GormDial = postgres.Open(dsn)
	case strings.HasPrefix(dsn, "sqlserver://"):
		GormDial = sqlserver.Open(dsn)
	default:
		return nil, errors.New("env DB_URL err")
	}
	return GormDial, nil
}
