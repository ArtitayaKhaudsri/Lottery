package service

import (
	"context"
	"encoding/json"
	"fmt"
	"lottery/config"
	"lottery/models"
	"math/rand/v2"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LotteryService struct {
	DB    *gorm.DB
	Redis *redis.Client
	cfg   *config.Config
}

func NewLotteryService(db *gorm.DB, redis *redis.Client, cfg *config.Config) *LotteryService {
	return &LotteryService{DB: db, Redis: redis, cfg: cfg}
}

func (h *LotteryService) Add(date string) {
	drawDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("invalid date")
		return
	}

	exist, err := h.Redis.Get(context.Background(), date).Result()
	if err == nil {
		fmt.Println(err)
		return
	}
	if exist != "" {
		fmt.Println("already created prize")
		return
	}

	var draw models.Draw
	if h.DB.Where("draw_date = ?", drawDate).First(&draw).Error == nil {
		fmt.Println("draw already exist")
		var oldPrize models.Prize
		if h.DB.Where("draw_id = ?", draw.DrawID).First(&oldPrize).Error == nil {
			fmt.Println("draw already exist")
			return
		}
	} else {
		draw = models.Draw{
			DrawDate: drawDate,
		}
		if err := h.DB.Create(&draw).Error; err != nil {
			fmt.Println("Date wrong:", err)
			return
		}
	}

	numbers := h.randomNumber()
	redisData := make([]models.RedisPrize, 0)
	for round, number := range numbers {
		prizeType := h.cfg.Type[round]
		Prize := models.Prize{
			PrizeType:     prizeType,
			WinningNumber: number,
			DrawID:        draw.DrawID,
		}
		if err := h.DB.Create(&Prize).Error; err != nil {
			fmt.Println("Something wrong:", err)
			return
		}
		redisData = append(redisData, models.RedisPrize{
			Number:    number,
			PrizeType: prizeType,
		})
	}
	fmt.Println("success")

	data, err := json.Marshal(redisData)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	if h.Redis.Set(context.Background(), date, data, time.Hour*24).Err() != nil {
		fmt.Println("redis set error:", err)
		return
	}
	fmt.Println("redis saved")
}

func (h *LotteryService) randomNumber() []string {
	result := make([]string, 0, 10)
	exists := make(map[string]bool)
	for len(result) < 10 {
		number := fmt.Sprintf("%06d", rand.IntN(1000000))
		if exists[number] {
			continue
		}
		exists[number] = true
		result = append(result, number)
	}
	return result
}
