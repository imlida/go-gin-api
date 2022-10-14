package mysql

import (
	"fmt"
	"time"

	"github.com/imlida/go-gin-api/configs"
	"github.com/imlida/go-gin-api/pkg/errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Predicate is a string that acts as a condition in the where clause
type Predicate string

var (
	EqualPredicate              = Predicate("=")
	NotEqualPredicate           = Predicate("<>")
	GreaterThanPredicate        = Predicate(">")
	GreaterThanOrEqualPredicate = Predicate(">=")
	SmallerThanPredicate        = Predicate("<")
	SmallerThanOrEqualPredicate = Predicate("<=")
	LikePredicate               = Predicate("LIKE")
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()
	GetDb(name string) *gorm.DB
	DbClose(db *gorm.DB) error
	GetDbs() map[string]*gorm.DB
}

type dbRepo struct {
	Dbs map[string]*gorm.DB
}

func New() (Repo, error) {
	cfg := configs.Get().MySQL

	var dbs = make(map[string]*gorm.DB)
	var db *gorm.DB
	var err error

	//遍历cfg,连接数据库
	for k, v := range cfg {
		db, err = dbConnect(v.User, v.Pass, v.Addr, v.Name, v.MaxOpenConn, v.MaxIdleConn, v.ConnMaxLifeTime)
		if err != nil {
			return nil, err
		}
		dbs[k] = db
	}

	return &dbRepo{
		Dbs: dbs,
	}, nil
}

func (d *dbRepo) i() {}

func (d *dbRepo) GetDb(name string) *gorm.DB {
	return d.Dbs[name]
}

func (d *dbRepo) DbClose(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *dbRepo) GetDbs() map[string]*gorm.DB {
	return d.Dbs
}

func dbConnect(user, pass, addr, dbName string, maxOpenConn, maxIdleConn int, connMaxLifeTime time.Duration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		user,
		pass,
		addr,
		dbName,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info), // 日志配置
	})

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("[db connection failed] Database name: %s", dbName))
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(maxOpenConn)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(maxIdleConn)

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(time.Minute * connMaxLifeTime)

	// 使用插件
	db.Use(&TracePlugin{})

	return db, nil
}
