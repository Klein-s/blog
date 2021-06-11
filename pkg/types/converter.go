package types

import "strconv"

//Int64ToString 将 int64 转为 string
func Int64ToString(num int64) string  {
	return strconv.FormatInt(num, 10)
}