package util

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

// Shuffle 打乱切片元素顺序
func Shuffle(src []int64) {
	rand.New(rand.NewSource(time.Now().Unix()))
	rand.Shuffle(len(src), func(i, j int) {
		src[i], src[j] = src[j], src[i]
	})
}

func GetNumbers(str string) []int64 {
	res := make([]int64, 0)
	re := regexp.MustCompile(`\d+`)

	// FindAllString 查找所有匹配的字符串
	matches := re.FindAllString(str, -1)

	// 打印所有匹配的数字
	for _, str := range matches {
		number, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			continue
		}
		res = append(res, number)
	}

	return res
}
