package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	seededRand *rand.Rand
	randOnce   sync.Once
)

func initRand() {
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max.
func RandomInt(min, max int) int {
	randOnce.Do(initRand)
	return min + seededRand.Intn(max-min+1)
}

// ShuffleStringSlice shuffles a slice of strings.
func ShuffleStringSlice(slice []string) []string {
	randOnce.Do(initRand)

	shuffled := make([]string, len(slice))
	copy(shuffled, slice)

	for i := range shuffled {
		j := seededRand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
