package utils

import "math/rand"

// RandomNum 随机数生成器
func RandomNum(s, e int) int {
	return s + rand.Intn(e)
}

// UColor make 0-255 number
func UColor() uint8 {
	return uint8(RandomNum(0, 255))
}
