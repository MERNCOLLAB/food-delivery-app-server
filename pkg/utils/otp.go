package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"time"

	"github.com/redis/go-redis/v9"
)

type OAuthTempData struct {
	Info      interface{}
	ExpiresAt time.Time
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

// OAuth Retrieved Data
func SetOAuthTemp(rdb *redis.Client, stateID string, data OAuthTempData, expiration time.Duration) error {
	ctx := context.Background()
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, "oauth"+stateID, b, expiration).Err()
}

func GetOAuthTemp(rdb *redis.Client, stateID string) (OAuthTempData, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, "oauth"+stateID).Result()
	if err != nil {
		return OAuthTempData{}, err
	}

	var data OAuthTempData
	err = json.Unmarshal([]byte(val), &data)
	return data, err
}

func DeleteOAuthTemp(rdb *redis.Client, stateID string) error {
	ctx := context.Background()
	return rdb.Del(ctx, "oauth:"+stateID).Err()
}

func GenerateStateID(rdb *redis.Client, info interface{}) string {
	stateID := GenerateUUIDStr()
	data := OAuthTempData{
		Info:      info,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	_ = SetOAuthTemp(rdb, stateID, data, 5*time.Minute)
	return stateID
}
