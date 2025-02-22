package oss

import "time"

type FileUpload interface {
	UploadFile(filename, storageDir string) (storagePath string, err error)
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
