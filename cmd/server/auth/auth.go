package auth

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"time"

	"github.com/tywkeene/go-tracker/cmd/server/db"
)

type RegisterAuth struct {
	Str             string
	Used            bool
	Timestamp       int64
	ExpireTimestamp int64
}

const authStrLen = 16

func Init(num int, expire time.Duration) error {
	authCount, err := db.GetRegisterAuthCount()
	if err != nil {
		return err
	}
	if authCount == 0 {
		log.Infof("Generating %d registration authorizations", num)
		for i := 0; i < num; i++ {
			auth := NewRegisterAuth(expire)
			err := db.InsertRegisterAuth(auth.Str, auth.Used, auth.Timestamp, auth.ExpireTimestamp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewRegisterAuth(expireDur time.Duration) *RegisterAuth {
	timestamp := time.Now().Unix()
	expire := (int64(expireDur.Seconds()) + timestamp)
	auth := &RegisterAuth{
		Str:             generateAuthStr(),
		Timestamp:       timestamp,
		ExpireTimestamp: expire,
		Used:            false,
	}
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

func ValidateRegisterAuth(authStr string) error {
	valid, err := db.IsAuthValid(authStr)
	if err != nil {
		return err
	}
	if valid == false {
		return fmt.Errorf("Invalid device register authorization: invalid auth string")
	}
	return db.SetAuthUsed(authStr, true)
}
