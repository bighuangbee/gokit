package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Options struct {
	Address  string
	UserName string
	Password string
	DBName   string
	Driver   string
	Logger   logger.Interface
	Charset  string
}

const (
	DbDriverMysql = "MySql"
)

func New(opt *Options) (*gorm.DB, error) {
	if opt.Driver == "" {
		opt.Driver = DbDriverMysql
	}
	if opt.Logger == nil {
		opt.Logger = logger.Default
	}

	var dialector gorm.Dialector

	if opt.Driver == DbDriverMysql {
		if opt.Charset == "" {
			opt.Charset = "utf8mb4"
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			opt.UserName, opt.Password, opt.Address, opt.DBName, opt.Charset)

		mysqlConfig := mysql.Config{
			DSN: dsn, // DSN data source name
			//DefaultStringSize:         191,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}
		dialector = mysql.New(mysqlConfig)
	} else {
		return nil, errors.New("DB Driver not found")
	}

	if db, err := gorm.Open(dialector); err != nil {
		return nil, err
	} else {
		db.Logger = opt.Logger
		db.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		}

		instance, err := db.DB()
		if err != nil {
			return nil, err
		}
		instance.SetMaxIdleConns(5)
		instance.SetMaxOpenConns(50)
		return db, nil
	}
}
