package util

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

// DeleteInt64 删除切片中所有等于val的元素
func DeleteInt64(slice []int64, val int64) []int64 {
	for i, item := range slice {
		if item == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	// 如果没有找到值，返回原始切片
	return slice
}

// Shuffle 打乱切片元素顺序
func Shuffle(src []int64) {
	rand.New(rand.NewSource(time.Now().Unix()))
	rand.Shuffle(len(src), func(i, j int) {
		src[i], src[j] = src[j], src[i]
	})
}

// GetNumbers 从字符串中解析出多个数字
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
