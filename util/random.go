package util

import (
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomVerificationCode(n int) []byte {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := strconv.Itoa(RandomInt(0, 9))
		sb.WriteString(c)
	}
	return []byte(sb.String())
}

func RandomUUID(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(uuid.NewString())
	}
	return sb.String()
}
