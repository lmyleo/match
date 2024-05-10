package util

import (
	"math/rand"
	"time"
)

func Shuffle(src []int64) {
	rand.New(rand.NewSource(time.Now().Unix()))
	rand.Shuffle(len(src), func(i, j int) {
		src[i], src[j] = src[j], src[i]
	})
}
