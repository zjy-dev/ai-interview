package data

import (
	"ai-interview/internal/conf"

	"database/sql"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewInterviewRepo)

// Data 封装数据库和缓存客户端
type Data struct {
	db  *sql.DB
	rdb *redis.Client
}

// NewData 初始化数据层
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)

	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}
	helper.Info("connected to mysql")

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})

	cleanup := func() {
		helper.Info("closing the data resources")
		_ = db.Close()
		_ = rdb.Close()
	}

	return &Data{db: db, rdb: rdb}, cleanup, nil
}
