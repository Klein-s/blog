package types

import (
	logger2 "goblog/pkg/logger"
	"strconv"
)

//Int64ToString 将 int64 转为 string
func Int64ToString(num int64) string  {
	return strconv.FormatInt(num, 10)
}

// Uint64ToString 将 Uint64 转为 string
func Uint64ToString(num uint64) string  {
	return  strconv.FormatUint(num, 10)
}

//StringToInt 将字符串转为 int
func StringToInt(str string) int  {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger2.LogError(err)
	}
	return i
}