package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"time"

	"github.com/redis/go-redis/v9"
)

type UserTempData struct {
	Info interface{}
}

func GenerateOTP() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

// OTP only
func SetOTP(rdb *redis.Client, phone, otp string, expiration time.Duration) error {
	ctx := context.Background()
	return rdb.Set(ctx, "otp:"+phone, otp, expiration).Err()
}

func GetOTP(rdb *redis.Client, phone string) (string, error) {
	ctx := context.Background()
	return rdb.Get(ctx, "otp:"+phone).Result()
}

func DeleteOTP(rdb *redis.Client, phone string) error {
	ctx := context.Background()
	return rdb.Del(ctx, "otp:"+phone).Err()
}

// Temporary User Data
func SetTempUser(rdb *redis.Client, redisKey string, data UserTempData, expiration time.Duration) error {
	ctx := context.Background()
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, "oauth:"+redisKey, b, expiration).Err()
}

func GetTempUser(rdb *redis.Client, redisKey string) (UserTempData, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, "oauth:"+redisKey).Result()
	if err != nil {
		return UserTempData{}, err
	}

	var data UserTempData
	err = json.Unmarshal([]byte(val), &data)
	return data, err
}

func DeleteTempUser(rdb *redis.Client, redisKey string) error {
	ctx := context.Background()
	return rdb.Del(ctx, "oauth:"+redisKey).Err()
}

func SetTempCustomer(rdb *redis.Client, info interface{}) string {
	redisKey := GenerateUUIDStr()
	data := UserTempData{
		Info: info,
	}

	_ = SetTempUser(rdb, redisKey, data, 20*time.Minute)
	return redisKey
}
