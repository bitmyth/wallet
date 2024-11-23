package factory

import (
	"errors"
	"github.com/bitmyth/walletserivce/config"
	"github.com/bitmyth/walletserivce/db"
	"github.com/bitmyth/walletserivce/wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Factory interface {
	Config() *config.Config
	DB() (*db.DB, error)
	Redis() (*db.Redis, error)
	Logger() *zap.SugaredLogger
	WalletController() *wallet.Controller
	RegisterRoutes(router *gin.Engine)
}

type Default struct {
	logger           *zap.SugaredLogger
	config           *config.Config
	db               *db.DB
	redis            *db.Redis
	walletController *wallet.Controller
}

func (d *Default) RegisterRoutes(router *gin.Engine) {
	d.WalletController().RegisterRoutes(router)
}

func (d *Default) WalletController() *wallet.Controller {
	if d.walletController == nil {
		d.walletController = wallet.NewController(d)
	}
	return d.walletController
}

func New() (Factory, error) {
	f := &Default{
		logger: logger(),
	}
	c, err := config.NewConfig()
	if err != nil {
		f.Logger().Error(err)
		return nil, err
	}
	f.config = c

	return f, nil
}

func (d *Default) Config() *config.Config {
	return d.config
}

func (d *Default) DB() (*db.DB, error) {
	conf := d.Config()
	var err error
	if d.db != nil {
		if d.db.IsAlive() {
			return d.db, nil
		}
		d.db = nil
	}

	d.db, err = db.Open(conf)
	return d.db, err
}

func (d *Default) Redis() (*db.Redis, error) {
	conf := d.Config()
	var err error
	if d.redis != nil {
		if d.redis.IsAlive() {
			return d.redis, nil
		}
		d.redis = nil
	}

	d.redis, err = db.OpenRedis(conf)
	return d.redis, err
}

func (d *Default) Logger() *zap.SugaredLogger {
	return d.logger
}

func logger() *zap.SugaredLogger {
	z, _ := zap.NewDevelopment(zap.AddCallerSkip(1))
	l := z.Sugar()
	return l
}

type TestingFactory struct {
	logger           *zap.SugaredLogger
	config           *config.Config
	walletController *wallet.Controller
}

func (t TestingFactory) RegisterRoutes(router *gin.Engine) {
	t.WalletController().RegisterRoutes(router)
}

func (t TestingFactory) WalletController() *wallet.Controller {
	if t.walletController == nil {
		t.walletController = wallet.NewController(t)
	}
	return t.walletController
}

func NewTesting() (Factory, error) {
	f := &TestingFactory{
		logger: logger(),
	}
	c, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	f.config = c

	return f, nil
}

func (t TestingFactory) Config() *config.Config {
	return t.config
}

func (t TestingFactory) DB() (*db.DB, error) {
	return nil, errors.New("db failed")
}

func (t TestingFactory) Redis() (*db.Redis, error) {
	return nil, errors.New("redis failed")
}

func (t TestingFactory) Logger() *zap.SugaredLogger {
	return t.logger
}
