/* 生成一个[0, n)的随机数 */
package utils

import (
	"math/rand"
	"time"
)

func RandInt(n int) int {
	if n <= 0 {
		return 0
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	return random.Intn(n)
}
