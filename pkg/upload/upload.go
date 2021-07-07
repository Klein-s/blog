package upload

import (
	"encoding/json"
	"fmt"
	logger2 "goblog/pkg/logger"
	"goblog/pkg/types"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"time"
)

type Image struct {
	Uploaded bool `json:"uploaded"`
	Url string `json:"url"`
}

/**
	上传图片
 	dir 为uploads 下级目录 如 images
	dirType 最后一级目录 如 articles
 */
func Upload(file multipart.File, dir string, dirType string)  []byte{
	dir = "./uploads/" + dir + "/" + dirType
	fileExist := checkFileIsExist(dir)
	fmt.Println("check file 1-----------------")
	fmt.Println("-------------------------fileExist:", fileExist)
	filename := types.Int64ToString(time.Now().Unix())+ types.Int64ToString(rand.Int63n(9999)) + ".png"
	fileDir := dir + "/" + filename
	if fileExist { //如果文件夹存在
		//f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件

		f, err := os.OpenFile(fileDir, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger2.LogError(err)
			return nil
		}
		defer f.Close()
		io.Copy(f, file)
		fmt.Println("文件夹存在")
	} else { //不存在文件夹时 先创建文件夹再上传
		err1 := os.MkdirAll(dir, os.ModePerm) //创建文件夹
		if err1 != nil {
			logger2.LogError(err1)
			return nil
		}

		fmt.Println("文件夹不存在")
		fmt.Println("文件夹创建成功！")
		f, err := os.OpenFile(fileDir, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger2.LogError(err)
			return nil
		}
		defer f.Close()
		io.Copy(f, file)
	}
	data := &Image{
		Uploaded: true,
		Url: "/uploads/images/articles/" + filename,
	}
	msg, _ := json.Marshal(data)

	return msg
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