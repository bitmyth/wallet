package route

import (
	"github.com/bitmyth/walletserivce/factory"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func checkDB(f factory.Factory) func(c *gin.Context) {
	return func(c *gin.Context) {
		if _, err := f.DB(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func checkRedis(f factory.Factory) func(c *gin.Context) {
	return func(c *gin.Context) {
		if _, err := f.Redis(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redis connection failed"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Router(f factory.Factory) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// Logging to a file.
	logFile, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(logFile)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	router.Use(checkDB(f), checkRedis(f))

	return router
}
