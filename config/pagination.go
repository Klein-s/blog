package config

import "goblog/pkg/config"

func init()  {
	config.Add("pagination", config.StrMap{
		//默认每页条数
		"perpage" : 10,

		//url 中分辨多少页的 参数

		"url_query":"page",
	})
}