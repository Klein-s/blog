package test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestAllPages(t *testing.T)  {

	baseUrl := "http://127.0.0.1:8005"

	//声明初始化测试数据

	var tests = []struct{
		method string
		url string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/3", 200},
		{"GET", "/articles/3/edit", 200},
		{"POST", "/articles/3", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/2/delete", 404},
	}

	//遍历所有测试
	for _, test := range tests{
		t.Logf("当前请求 URl: %v \n", test.url)
		var (
			resp *http.Response
			err error
		)
		//请求获取响应
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseUrl+test.url, data)
		default:
			resp, err = http.Get(baseUrl + test.url)
		}

		//断言
		assert.NoError(t, err, "请求" + test.url + "时报错")
		assert.Equal(t, test.expected, resp.StatusCode, test.url +
			"应返回状态码" + strconv.Itoa(test.expected))
	}
}