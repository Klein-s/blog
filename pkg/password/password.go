package password

import (
	logger2 "goblog/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

//hash 使用 bcrypt 对密码进行加密
func Hash(password string) string  {
	// generateFromPassword 第二个参数是 cost 值， 建议大于12， 数值越大耗时越长
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger2.LogError(err)

	return string(bytes)
}

//checkHash 对吧明文密码和数据库的哈希值

func CheckHas(password, hash string) bool  {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	logger2.LogError(err)

	return err == nil
}

// isHasded 判断字符串是否是哈希过的数据
func IsHashed(str string) bool  {
	//加密后的长度等于 60
	return len(str) == 60
}