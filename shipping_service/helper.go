package main

import (
	"fmt"
	"math/rand"
	"time"
)

func CreateTrackingId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%c%c-%s-%s",
		generateRandomCharacter(),
		generateRandomCharacter(),
		generateRandomNumber(5),
		generateRandomNumber(9),
	)
}

func generateRandomCharacter() uint32 {
	return 65 + uint32(rand.Intn(25))
}

func generateRandomNumber(length int) string {
	s := ""
	for i := 0; i < length; i++ {
		s = fmt.Sprintf("%s%d", s, rand.Intn(10))
	}

	return s
}
