package tools

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

//获取文件名，不包含后缀名
func GetFilename(filePath string)string{
	filePath = path.Base(filePath)
	ext := path.Ext(filePath)
	return filePath[0:len(filePath) - len(ext)]
}

//目录是否存在
func PathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//获取当前目录的名称
func GetPathName()string{
	s, _ := os.Getwd()
	return path.Base(s)
}

func GetFilesByPath(basePath string) ([]string) {
	fileArr := []string{}
	fs,_:= ioutil.ReadDir(basePath)
	for _,file := range fs{
		if file.IsDir(){
			fileArr = append(fileArr, GetFilesByPath(basePath+file.Name()+"/")...)
		}else{
			fileArr = append(fileArr, basePath + "/" + file.Name())
		}
	}
	return fileArr
}


func SaveByFileHeader(header *multipart.FileHeader, filename string)error{
	file, err := header.Open()
	if err != nil{
		return err
	}

	os.MkdirAll(filepath.Dir(filename), 0755)
	f1,_ := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0755)
	defer f1.Close()

	for{
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF{
			return err
		}

		if n == 0{
			break
		}
		if _, err = f1.Write(buf); err != nil{
			return err
		}
	}
	return nil
}
