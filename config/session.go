package config

import "goblog/pkg/config"

func init()  {
	config.Add("session", config.StrMap{

		//目前只支持 cookie

		"default":config.Env("SESSION_DRIVER", "cookie"),

		//会话的 cookie 名称
		"session_name":config.Env("SESSION_NAME", "blog-session"),
	})
}
