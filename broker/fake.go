package broker

import (
	"math/rand"
	"time"
)

// there are mostly stub function which will be removed later
// for db
func getUserKeys(user string) []string {
	return nil
}

// for db
func getParent(key string) (string, error) {
	return "", nil
}

// for session/cookie manager
func getUserFromCookie(user string) string {
	return ""
}

// just small experiment, temp probably

func generateId() int64 {

	now := time.Now()
	tnano := now.UnixNano()

	randint := rand.Int31()

	return (tnano << 32) | int64(randint)
}
