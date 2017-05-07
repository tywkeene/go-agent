package auth

import (
	"fmt"
	"math/rand"
	"time"
)

type RegisterAuth struct {
	Str       string
	Timestamp int64
	Expire    time.Duration
	Used      bool
}

var validRegisterAuths []*RegisterAuth

const authStrLen = 16

func Init() {
	validRegisterAuths = make([]*RegisterAuth, 0)
}

func NewRegisterAuth(expire time.Duration) *RegisterAuth {
	auth := &RegisterAuth{
		Str:       generateAuthStr(),
		Timestamp: time.Now().UnixNano(),
		Expire:    expire,
		Used:      false,
	}
	validRegisterAuths = append(validRegisterAuths, auth)
	return auth
}

func generateAuthStr() string {
	var random = []rune("ABCDEF1234567890")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, authStrLen)
	for i := range b {
		b[i] = random[rand.Intn(len(random))]
	}
	return string(b)
}

func InvalidateAuth(str string) error {
	for _, auth := range validRegisterAuths {
		if auth.Str == str {
			auth.Used = true
			return nil
		}
	}
	return fmt.Errorf("Could not find authorization to invalidate")
}

func ValidateRegisterAuth(authStr string) error {
	for _, auth := range validRegisterAuths {
		if auth.Str == authStr {
			if auth.Used == true {
				return fmt.Errorf("Invalid device register authorization: already used")
			}
			timeSince := time.Unix(0, auth.Timestamp)
			diff := time.Now().Sub(timeSince)
			if diff > auth.Expire {
				return fmt.Errorf("Invalid device register authorization: expired")
			}
			return nil
		}
	}
	return fmt.Errorf("Invalid device register authorization: invalid auth string")
}
