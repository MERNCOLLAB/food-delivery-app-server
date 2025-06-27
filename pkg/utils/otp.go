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
	M map[string]string
}{M: make(map[string]string)}

var OAuthTempStore = struct {
	sync.RWMutex
	M map[string]OAuthTempData
}{M: make(map[string]OAuthTempData)}

func GenerateOTP() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

func GenerateStateID(info interface{}) string {
	stateID := GenerateUUIDStr()
	OAuthTempStore.Lock()
	OAuthTempStore.M[stateID] = OAuthTempData{
		Info:      info,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	OAuthTempStore.Unlock()
	return stateID
}

func CleanMemory(phone, stateID string) {
	OtpStore.Lock()
	delete(OtpStore.M, phone)
	OtpStore.Unlock()

	OAuthTempStore.Lock()
	delete(OAuthTempStore.M, stateID)
	OAuthTempStore.Unlock()
}
