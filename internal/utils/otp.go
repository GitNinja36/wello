package utils

import (
	"sync"
	"time"
)

var otpStore = make(map[string]otpEntry)
var otpMutex sync.Mutex

type otpEntry struct {
	OTP     string
	Expires time.Time
}

func SaveOTP(key, otp string, ttl time.Duration) {
	otpMutex.Lock()
	defer otpMutex.Unlock()
	otpStore[key] = otpEntry{
		OTP:     otp,
		Expires: time.Now().Add(ttl),
	}
}

func VerifyOTP(key, otp string) bool {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	entry, exists := otpStore[key]
	if !exists || entry.OTP != otp {
		return false
	}
	if time.Now().After(entry.Expires) {
		delete(otpStore, key)
		return false
	}
	delete(otpStore, key)
	return true
}
