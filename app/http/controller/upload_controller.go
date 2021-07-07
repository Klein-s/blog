package controller

import (
	"encoding/json"
	"fmt"
	logger2 "goblog/pkg/logger"
	"goblog/pkg/types"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type UploadController struct {
	BaseController
}

type UploadImage struct {
	Uploaded bool `json:"uploaded"`
	Url string `json:"url"`
}
/**
图片上传
*/
func (*UploadController) UploadImage(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("upload")
	if err != nil {
		logger2.LogError(err)
		return
	}
	defer file.Close()


	fileExist := checkFileIsExist("./uploads")
	fmt.Println("check file 1-----------------")
	fmt.Println("-------------------------fileExist:", fileExist)
	filename := types.Int64ToString(time.Now().Unix())+ types.Int64ToString(rand.Int63n(9999)) + ".png"
	fileDir := "./uploads/" + filename
	if fileExist { //如果文件夹存在
		//f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件

		f, err := os.OpenFile(fileDir, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger2.LogError(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		fmt.Println("文件夹存在")
	} else { //不存在文件夹时 先创建文件夹再上传
		err1 := os.Mkdir("./uploads", os.ModePerm) //创建文件夹
		if err1 != nil {
			logger2.LogError(err1)
			return
		}

		fmt.Println("文件夹不存在")
		fmt.Println("文件夹创建成功！")
		f, err := os.OpenFile(fileDir, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger2.LogError(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
	data := &UploadImage{
		Uploaded: true,
		Url: "/uploads/" + filename,
	}
	msg, _ := json.Marshal(data)
	w.Header().Set("content-type","text/json")
	w.Write(msg)
}

/**
显示图片
 */
func (*UploadController) ShowImages(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open("."+ r.URL.Path)
	if err != nil {
		logger2.LogError(err)
	}
	defer file.Close()
	buff, err := ioutil.ReadAll(file)

	w.Write(buff)

}
//检查目录是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Print(filename + " not exist")
		exist = false
	}
	return exist
}
