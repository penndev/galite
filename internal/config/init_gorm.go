package config

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/penndev/galite/internal/lib"
	loggerNal "github.com/penndev/galite/internal/logger"
	loggerPkg "github.com/penndev/galite/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBMaxOpenConns = 100

// 数据库连接实例，
// 将多个数据库连接进行抽象，减少不同数据库配置依赖
// 不论什么数据库最终提供给gram的都为 gorm.Dialector
func gormDial(dsn string) (gorm.Dialector, error) {
	var GormDial gorm.Dialector
	switch {
	case strings.HasPrefix(dsn, "mariadb://"):
		GormDial = mysql.Open(strings.TrimPrefix(dsn, "mariadb://"))
	case strings.HasPrefix(dsn, "mysql://"):
		GormDial = mysql.Open(strings.TrimPrefix(dsn, "mysql://"))
	case strings.HasPrefix(dsn, "sqlite://"):
		GormDial = sqlite.Open(strings.TrimPrefix(dsn, "sqlite://"))
		// https://github.com/glebarez/sqlite/issues/52
		DBMaxOpenConns = 1
	case strings.HasPrefix(dsn, "postgres://"):
		GormDial = postgres.Open(dsn)
	case strings.HasPrefix(dsn, "sqlserver://"):
		GormDial = sqlserver.Open(dsn)
	default:
		return nil, errors.New("env DB_URL err")
	}
	return GormDial, nil
}

// 配置数据库连接信息
func gormInit(dialector gorm.Dialector, Logger logger.Interface) error {
	// var err error
	// var dataBase *gorm.DB
	dataBase, err := gorm.Open(dialector, &gorm.Config{
		Logger: Logger, // 重写日志
		// DisableForeignKeyConstraintWhenMigrating: true, // 禁止物理外键约束
	})
	if err != nil {
		return err
	}

	sqlDB, err := dataBase.DB()
	if err != nil {
		return err
	}

	// 连接池，程序关闭时回收
	closeList = append(closeList, sqlDB)
	// 配置连接池参数
	if DBMaxOpenConns > 0 {
		// 最大连接数
		sqlDB.SetMaxOpenConns(DBMaxOpenConns)
		// 最大空闲数
		sqlDB.SetMaxIdleConns(DBMaxOpenConns / 5)
		// 最大存活时间
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	lib.SetGorm(dataBase)
	return nil
}

// 设置gorm数据库的实例
func InitGorm() error {
	gormDial, err := gormDial(cfg.Database.Dsn)
	if err != nil {
		return err
	}
	if cfg.Mode == ModeDev {
		gormLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  cfg.Logger.Level.GormLevel(),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true, //颜色控制
		})
		if err := gormInit(gormDial, gormLogger); err != nil { // 初始化
			return err
		}
	} else {
		gormLogger := loggerPkg.ZapGormLogger(
			loggerNal.GormZapLogger,
			200*time.Millisecond,
		)
		if err := gormInit(gormDial, gormLogger); err != nil { // 初始化
			return err
		}
	}
	return nil
}
