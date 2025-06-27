package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type OAuthTempData struct {
	Info      interface{}
	ExpiresAt time.Time
}

var OtpStore = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var OAuthTempStore = struct {
	sync.RWMutex
	m map[string]OAuthTempData
}{m: make(map[string]OAuthTempData)}

func GenerateOTP() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

func GenerateStateID(info interface{}) string {
	stateID := GenerateUUIDStr()
	OAuthTempStore.Lock()
	OAuthTempStore.m[stateID] = OAuthTempData{
		Info:      info,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	OAuthTempStore.Unlock()
	return stateID
}
