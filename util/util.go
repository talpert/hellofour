package util

import (
	"math/rand"
	"time"

	"github.com/satori/go.uuid"
)

const (
	RAND_LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandString(n int, t time.Time) string {
	src := rand.NewSource(t.UTC().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = RAND_LETTERS[src.Int63()%int64(len(RAND_LETTERS))]
	}
	return string(b)
}

// a random uuid to use as a namespace
var nsUUID = uuid.Must(uuid.FromString("6159f029-a176-4177-9251-01515315c662"))

// generate a UUID by iterating over the strategies, without throwing an error
func GenerateUUID() uuid.UUID {
	/*****
	//eventually do this once we upgrade to new lib versions.
	// for now, still use this function as an unchanging interface,
	// but it does not do anything
	****/

	id, err := uuid.NewV1()
	if err != nil {
		id, err = uuid.NewV4()
		if err != nil {
			return uuid.NewV5(nsUUID, time.Now().String())
		}
	}

	return id
}
