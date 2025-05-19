package oss

import (
	"io"
	"time"
)

type FileUpload interface {
	UploadFile(filename, storageFilename string) (err error)
	UploadReader(reader io.Reader, storageFilename string) (err error)
	DeleteDir(prefix string) (err error)
	Sign() (token string, err error)
}

type FileDownload interface {
	Download(filepath, localPath string) error
}

type Oss struct {
	FileDownload FileDownload
	FileUpload   FileUpload
}

func New(accessKey, secretKey, bucket string, tokenExpire time.Duration) *Oss {
	ossInstance := NewQiniu(accessKey, secretKey, bucket, time.Minute*tokenExpire)
	return &Oss{
		FileDownload: ossInstance,
		FileUpload:   ossInstance,
	}
}
