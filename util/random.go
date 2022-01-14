package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/teris-io/shortid"
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

func RandomUnique(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		gen, _ := shortid.Generate()
		sb.WriteString(gen)
	}
	return sb.String()
}
