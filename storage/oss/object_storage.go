package oss

type FileUpload interface {
	UploadFile(filename, storageDir string) error
	Sign() (token string, err error)
}

type FileDownload interface {
	Download(filepath, localPath string) error
}
