package oss

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/downloader"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"path/filepath"
	"time"
)

type Qiniu struct {
	accessKey            string
	secretKey            string
	bucket               string
	expire               time.Duration
	uploadManager        *uploader.UploadManager
	token                uptoken.Provider
	downloadManager      *downloader.DownloadManager
	downloadURLsProvider downloader.DownloadURLsProvider
}

func NewQiniu(accessKey string, secretKey string, bucket string, expire time.Duration) *Qiniu {
	mac := credentials.NewCredentials(accessKey, secretKey)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})

	downloadManager := downloader.NewDownloadManager(&downloader.DownloadManagerOptions{})
	downloadURLsProvider := downloader.SignURLsProvider(downloader.NewDefaultSrcURLsProvider(mac.AccessKey, nil), downloader.NewCredentialsSigner(mac), nil)

	qiniu := &Qiniu{accessKey: accessKey, secretKey: secretKey, bucket: bucket, expire: expire,
		uploadManager:        uploadManager,
		downloadManager:      downloadManager,
		downloadURLsProvider: downloadURLsProvider,
	}

	qiniu.refreshToken()
	ticker := time.NewTicker(expire)
	go func() {
		for {
			select {
			case <-ticker.C:
				qiniu.refreshToken()
			}
		}
	}()

	return qiniu
}

func (q *Qiniu) UploadFile(filename, storageDir string) (storagePath string, err error) {
	name := filepath.Join(storageDir, filepath.Base(filename))

	err = q.uploadManager.UploadFile(context.Background(), filename, &uploader.ObjectOptions{
		BucketName: q.bucket,
		ObjectName: &name,
		FileName:   name,
	}, nil)

	return name, err
}

func (q *Qiniu) Sign() (token string, err error) {
	return q.token.GetUpToken(context.Background())
}

func (q *Qiniu) refreshToken() error {
	mac := credentials.NewCredentials(q.accessKey, q.secretKey)
	putPolicy, err := uptoken.NewPutPolicy(q.bucket, time.Now().Add(q.expire))
	if err != nil {
		return err
	}

	newToken := uptoken.NewSigner(putPolicy, mac)
	q.token = newToken
	return nil
}

func (q *Qiniu) Download(filepath, localPath string) error {
	_, err := q.downloadManager.DownloadToFile(context.Background(), filepath, localPath, &downloader.ObjectOptions{
		GenerateOptions:      downloader.GenerateOptions{BucketName: q.bucket},
		DownloadURLsProvider: q.downloadURLsProvider,
	})
	if err != nil {
		return err
	}

	return err
}
