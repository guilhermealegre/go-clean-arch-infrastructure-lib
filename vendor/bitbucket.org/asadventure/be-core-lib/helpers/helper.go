package helpers

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func CurrentDate() string {
	return time.Now().Format(time.RFC3339)
}

func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) &&
			((i+1 < length && unicode.IsLower(runes[i+1])) ||
				unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func RandStringBytes(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Percentage(x float64) float64 {
	return float64(int64(Round(x) * 100.0))
}

func Round(x float64) float64 {
	return float64(int64(x*100+0.5)) / 100
}

// LowerSlice lowers a slice of strings
func LowerSlice(items []string) []string {
	for i, item := range items {
		items[i] = TrimAndLowerStr(item)
	}

	return items
}

func IntInSlice(item int, list []int) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func TrimAndLowerStr(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

func DistinctArray(str []string) []string {
	processed := make(map[string]bool)
	var newStr []string

	for _, item := range str {
		if _, ok := processed[item]; !ok {
			newStr = append(newStr, item)
			processed[item] = true
		}
	}

	return newStr
}

func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func GetString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
