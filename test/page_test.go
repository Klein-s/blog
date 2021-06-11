package test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHomePage(t *testing.T)  {
	baseUrl := "http://127.0.0.1:8005"

	//请求  模拟用户访问浏览器
	var (
		resp *http.Response
		err error
	)
	resp, err = http.Get(baseUrl + "/")

	//检测  是否无措且 200
	assert.NoError(t, err, "有错误发生， err 不为空")
	assert.Equal(t, 200, resp.StatusCode, "返回状态码 200")
}