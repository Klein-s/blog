package requests

import (
	"errors"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"goblog/pkg/model"
	"strconv"
	"strings"
	"unicode/utf8"
)

func init() {

	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {

		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tableName := rng[0]
		dbFiled := rng[1]
		val := value.(string)

		var count int64

		model.DB.Table(tableName).Where(dbFiled + "= ?", val).Count(&count)

		if count != 0 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已被占用", val)
		}
		return nil
	})

	govalidator.AddCustomRule("max_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:"))
		if valLength > l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}
		return nil
	})

	govalidator.AddCustomRule("min_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:"))
		if valLength < l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度需大于 %d 个字", l)
		}
		return nil
	})

	govalidator.AddCustomRule("required_id", func(field string, rule string, message string, value interface{}) error {

		if value == 0 {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("分类为必填")
		}
		return nil
	})
}
