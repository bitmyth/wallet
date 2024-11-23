package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bitmyth/walletserivce/config"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DB struct {
	*sql.DB
}

func (d DB) IsAlive() bool {
	err := d.Ping()
	return err == nil
}

func Open(conf *config.Config) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname, conf.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

type Redis struct {
	*redis.Client
}

func (r Redis) IsAlive() bool {
	_, err := r.Ping(context.Background()).Result()
	return err == nil
}

func OpenRedis(c *config.Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		Client: rdb,
	}, nil
}
