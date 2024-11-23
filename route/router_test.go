package route

import (
	"github.com/bitmyth/walletserivce/config"
	"github.com/bitmyth/walletserivce/factory"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	config.SetConfigPath("../")
	m.Run()
}

func TestRouter(t *testing.T) {
	f, err := factory.New()
	if err != nil {
		return
	}

	r := Router(f)
	f.RegisterRoutes(r)
	t.Log("routes count", len(r.Routes()))
	if len(r.Routes()) == 0 {
		t.Error("routes count is not 0")
	}
}

func TestCheckDBFailed(t *testing.T) {
	mockFactory, _ := factory.NewTesting()
	router := gin.New()
	router.Use(checkDB(mockFactory))
	w := flight(router)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "database connection failed"}`, w.Body.String())
}

func TestCheckDBSuccess(t *testing.T) {
	config.SetConfigPath("../")

	router := gin.New()

	mockFactory, _ := factory.New()
	router.Use(checkDB(mockFactory))

	w := flight(router)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "success"}`, w.Body.String())
}

func TestCheckRedisFailed(t *testing.T) {
	router := gin.New()

	mockFactory, _ := factory.NewTesting()
	router.Use(checkRedis(mockFactory))

	w := flight(router)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "redis connection failed"}`, w.Body.String())
}

func TestCheckRedisSuccess(t *testing.T) {
	config.SetConfigPath("../")

	router := gin.New()

	mockFactory, _ := factory.New()
	router.Use(checkRedis(mockFactory))

	w := flight(router)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "success"}`, w.Body.String())
}

func flight(router *gin.Engine) *httptest.ResponseRecorder {
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)
	return w
}
