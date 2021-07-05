package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/model/article"
)

// 文章表单验证
func ValidateArticleForm(data article.Article) map[string][]string  {

	//定制认证规则
	rules := govalidator.MapData{
		"title": []string{"required", "min_cn:3", "max_cn:40"},
		"body": []string{"required", "min_cn:10"},
		"category_id":[]string{"required"},
	}

	//定制错误消息
	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"min_cn:标题长度大于3",
			"max_cn:标题长度小于40",
		},
		"body":[]string{
			"required:文章内容为必填",
			"min_cn:长度需要大于10",
		},
		"category_id":[]string{
			"required:请选择分类",
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
