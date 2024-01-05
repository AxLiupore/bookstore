package service

import (
	"bookstore/conf"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

func UploadAvatarToLocalStatic(file multipart.File, id uint, username string) (filePath string, err error) {
	bId := strconv.Itoa(int(id)) // 路径拼接
	basePath := "." + conf.AvatarPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + username + ".jpg" // todo 把 file 的后缀提取出来
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(avatarPath, content, 0666)
	if err != nil {
		return
	}
	return "user" + bId + "/" + username + ".jpg", nil
}

// DirExistOrNot 判断文件夹路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return false
	}
	return true
}
