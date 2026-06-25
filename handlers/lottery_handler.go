package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"lottery/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LotteryHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewLotteryHandler(db *gorm.DB, redis *redis.Client) *LotteryHandler {
	return &LotteryHandler{DB: db, Redis: redis}
}

func (h *LotteryHandler) CheckNumber(c *gin.Context) {
	date := c.Query("date")
	number := c.Query("number")
	errormessage := validateRequest(date, number)
	if errormessage != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errormessage})
		return
	}
	exist, err := h.Redis.Get(context.Background(), date).Result()
	if err == nil {
		var prizes []models.RedisPrize
		message := "Sorry, you are not winner"
		if err := json.Unmarshal([]byte(exist), &prizes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redis data invalid"})
			return
		}
		for _, prize := range prizes {
			if prize.Number == number {
				message = fmt.Sprintf("Congratulations! You got %s prize", prize.PrizeType)
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": message})
		return
	}

	var draw models.Draw
	if err := h.DB.Where("draw_date = ?", date).First(&draw).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, draw not announced yet"})
		return
	}
	var prize models.Prize
	err = h.DB.Where("draw_id = ? and winning_number = ?", draw.DrawID, number).First(&prize).Error
	if err == nil {
		message := fmt.Sprintf("Congratulations! You got %s prize", prize.PrizeType)
		c.JSON(http.StatusOK, gin.H{"message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sorry, you are not winner"})
	return
}

func validateRequest(date string, number string) string {
	if date == "" || number == "" {
		return "date and number are required"
	}
	if len(number) != 6 {
		return "number must be 6 digit"
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return "invalid date format"
	}
	return ""
}

func (h *LotteryHandler) SearchNumber(c *gin.Context) {
	number := c.Query("number")
	var result models.PrizeResponse
	err := h.DB.Table("prizes").
		Select("draws.draw_date as date, prizes.winning_number as number, prizes.prize_type as type").
		Joins("JOIN draws ON draws.draw_id = prizes.draw_id").
		Where("prizes.winning_number = ?", number).
		Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
