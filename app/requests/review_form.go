package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/model/review"
)

func ValidateReviewForm(data review.Review) map[string][]string {

	rules := govalidator.MapData{
		"content": []string{"required"},
	}

	messages := govalidator.MapData{
		"content": []string{
			"required:评论内容必须填写",
		},
	}

	opts := govalidator.Options{
		Data: &data,
		Rules: rules,
		Messages: messages,
		TagIdentifier: "valid",
	}

	return govalidator.New(opts).ValidateStruct()
}