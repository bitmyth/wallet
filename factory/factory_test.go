package factory

import (
	"github.com/bitmyth/walletserivce/config"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestMain(m *testing.M) {
	config.SetConfigPath("../")
	m.Run()
}

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Error(err)
	}
}

func TestDefault_RegisterRoutes(t *testing.T) {
	factory, _ := New()
	router := gin.New()
	factory.RegisterRoutes(router)
	if len(router.Routes()) == 0 {
		t.Error("expect routes not 0")
	}
}

func TestDefault_DB(t *testing.T) {
	factory, _ := New()
	_, err := factory.DB()
	if err != nil {
		t.Error(err)
	}
}

func TestDefault_Redis(t *testing.T) {
	factory, _ := New()
	_, err := factory.Redis()
	if err != nil {
		t.Error(err)
	}
}
func TestDefault_Logger(t *testing.T) {
	factory, _ := New()
	l := factory.Logger()
	if l == nil {
		t.Error("logger is nil")
	}
}

func TestTestingFactory_Config(t *testing.T) {
	factory, _ := NewTesting()
	if factory.Config() == nil {
		t.Error("config is nil")
	}
}

func TestTestingFactory_DB(t *testing.T) {
	factory, _ := NewTesting()
	db, _ := factory.DB()

	if db != nil {
		t.Error("expect db is nil")
	}
}

func TestTestingFactory_Redis(t *testing.T) {
	factory, _ := NewTesting()
	db, _ := factory.Redis()

	if db != nil {
		t.Error("expect redis is nil")
	}
}

func TestTesting_RegisterRoutes(t *testing.T) {
	factory, _ := NewTesting()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	factory.RegisterRoutes(router)
	if len(router.Routes()) == 0 {
		t.Error("expect routes not 0")
	}
}
