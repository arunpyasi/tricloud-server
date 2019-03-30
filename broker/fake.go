package broker

import (
	"math/rand"
	"time"
)

// there are mostly stub function which will be removed later

func getUserKeys(user string) []string {
	return nil
}

func getParent(key string) (string, error) {
	return "", nil
}

// just small experiment, temp probably

func generateId() int64 {

	now := time.Now()
	tnano := now.UnixNano()

	randint := rand.Int31()

	return (tnano << 32) | int64(randint)
}
