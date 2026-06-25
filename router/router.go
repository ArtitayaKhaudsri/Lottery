package routers

import (
	handler "lottery/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, redis *redis.Client) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	lotteryHandler := handler.NewLotteryHandler(db, redis)
	lottery := r.Group("/lottery")
	{
		lottery.GET("/check", lotteryHandler.CheckNumber)
		lottery.GET("/search", lotteryHandler.SearchNumber)
	}

	return r
}
