package tracking

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateTrackingCode() string {
	timestamp := time.Now().Unix() % 1000
	random := rand.Intn(100000)
	return fmt.Sprintf("BR%03d%05d", timestamp, random)
}

func GenerateUniqueTrackingCode(existsFunc func(string) bool) string {
	for attempts := 0; attempts < 10; attempts++ {
		code := GenerateTrackingCode()
		if !existsFunc(code) {
			return code
		}
	}
	// Fallback com mais entropia
	timestamp := time.Now().UnixNano() % 100000000
	return fmt.Sprintf("BR%08d", timestamp)
}
