package controller

import (
	logger2 "goblog/pkg/logger"
	"goblog/pkg/upload"
	"io/ioutil"
	"net/http"
	"os"
)

type UploadController struct {
	BaseController
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
	dir := r.URL.Query().Get("dir")
	dirType := r.URL.Query().Get("type")
	msg := upload.Upload(file, dir, dirType)

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

