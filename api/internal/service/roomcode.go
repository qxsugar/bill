package service

import (
	"math/rand"
)

// hasPairedConsecutive 判断 4 位数字串是否存在「两两连号」，
// 即任意相邻两位相同（如 1123 中的 11、2255 中的 22、55）。
func hasPairedConsecutive(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == s[i+1] {
			return true
		}
	}
	return false
}

// gen4DigitPreferred 生成 4 位房间码候选：
// 优先尝试带「两两连号」的号码，多次失败后退回任意 4 位数。
// exists 用于判重（已占用返回 true）。返回空串表示放弃 4 位空间。
func gen4DigitPreferred(exists func(code string) bool) string {
	// 先在「含两两连号」的候选里随机尝试。
	for i := 0; i < 40; i++ {
		code := randDigits(4)
		if hasPairedConsecutive(code) && !exists(code) {
			return code
		}
	}
	// 退而求其次：任意未占用的 4 位数。
	for i := 0; i < 40; i++ {
		code := randDigits(4)
		if !exists(code) {
			return code
		}
	}
	return ""
}

// gen5Digit 生成未占用的 5 位房间码（4 位耗尽后使用）。
func gen5Digit(exists func(code string) bool) string {
	for i := 0; i < 80; i++ {
		code := randDigits(5)
		if !exists(code) {
			return code
		}
	}
	return ""
}

// randDigits 生成 n 位数字串，首位允许为 0（房间码按字符串存储）。
func randDigits(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rand.Intn(10))
	}
	return string(b)
}
