package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/model/article"
)

// 文章表单验证
func ValidateArticleForm(data article.Article) map[string][]string  {

	//定制认证规则
	rules := govalidator.MapData{
		"title": []string{"required", "min:3", "max:40"},
		"body": []string{"required", "min:10"},
	}

	//定制错误消息
	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"min:标题长度大于3",
			"max:标题长度小于40",
		},
		"body":[]string{
			"required:文章内容为必填",
			"min:长度需要大于10",
		},
	}

	//配置初始化
	opts := govalidator.Options{
		Data: &data,
		Rules: rules,
		TagIdentifier: "valid",
		Messages: messages,
	}

	//开始验证
	return govalidator.New(opts).ValidateStruct()
}
