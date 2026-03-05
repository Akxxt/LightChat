package utils

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() {
	// 从环境变量获取配置，未获取则用默认值
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	dbStr := getEnv("REDIS_DB", "0")
	db, _ := strconv.Atoi(dbStr)

	// 创建Redis客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("Redis连接失败: " + err.Error())
	}
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		InitRedis()
	}
	return redisClient
}

// GenerateID 生成简单的唯一ID（简化版，实际可用UUID）
func GenerateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// getEnv 获取环境变量，不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}