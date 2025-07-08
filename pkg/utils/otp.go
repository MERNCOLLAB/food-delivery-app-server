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
func SetTempUser(rdb *redis.Client, stateID string, data UserTempData, expiration time.Duration) error {
	ctx := context.Background()
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, "oauth"+stateID, b, expiration).Err()
}

func GetTempUser(rdb *redis.Client, stateID string) (UserTempData, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, "oauth"+stateID).Result()
	if err != nil {
		return UserTempData{}, err
	}

	var data UserTempData
	err = json.Unmarshal([]byte(val), &data)
	return data, err
}

func DeleteTempUser(rdb *redis.Client, stateID string) error {
	ctx := context.Background()
	return rdb.Del(ctx, "oauth:"+stateID).Err()
}

func GenerateStateID(rdb *redis.Client, info interface{}) string {
	stateID := GenerateUUIDStr()
	data := UserTempData{
		Info: info,
	}

	_ = SetTempUser(rdb, stateID, data, 5*time.Minute)
	return stateID
}
