package wallet

import (
	"context"
	"github.com/bitmyth/walletserivce/config"
	"github.com/bitmyth/walletserivce/db"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
)

type dependency interface {
	Config() *config.Config
	DB() (*db.DB, error)
	Redis() (*db.Redis, error)
	Logger() *zap.SugaredLogger
	WalletController() *Controller
	RegisterRoutes(router *gin.Engine)
}

type Service struct {
	factory dependency
}

func NewService(factory dependency) *Service {
	return &Service{factory: factory}
}

func (s Service) GetBalance(ctx context.Context, username string) (float64, error) {
	logger := s.factory.Logger()

	d, err := s.factory.DB()
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	rdb, er := s.factory.Redis()
	if er != nil {
		logger.Error(err)
		return 0, err
	}

	var balance float64
	// check cache first
	b, err := rdb.Get(ctx, username).Result()
	if err == redis.Nil {
		// not in cache, get from DB
		var user User
		err = d.QueryRow("SELECT id, username, balance FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Balance)
		if err != nil {
			return 0, err
		}
		// cache the balance
		rdb.Set(ctx, username, user.Balance, 0)
		balance = user.Balance
	} else if err != nil {
		return 0, err
	} else {
		balance, _ = strconv.ParseFloat(b, 64)
	}

	return balance, nil
}
