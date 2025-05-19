package oss

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/downloader"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"io"
	"sync"
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
	bm                   *storage.BucketManager
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

	cfg := storage.Config{
		//Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	qiniu.bm = storage.NewBucketManager(mac, &cfg)
	return qiniu
}

func (q *Qiniu) UploadFile(filename, storageFilename string) (err error) {
	err = q.uploadManager.UploadFile(context.Background(), filename, &uploader.ObjectOptions{
		BucketName: q.bucket,
		ObjectName: &storageFilename,
		FileName:   storageFilename,
	}, nil)
	return err
}

func (q *Qiniu) UploadReader(reader io.Reader, storageFilename string) error {
	err := q.uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: q.bucket,
		ObjectName: &storageFilename,
		FileName:   storageFilename,
	}, nil)
	return err
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

// ======= delete =======

type DeleteResult struct {
	Total    int
	Success  int
	Failures map[string]error
	mu       sync.Mutex
}

func (q *Qiniu) DeleteDir(prefix string) error {
	keys, _ := listAllKeys(q.bm, q.bucket, prefix)
	if len(keys) == 0 {
		return nil
	}

	result := DeleteResult{
		Failures: make(map[string]error),
		Total:    len(keys),
	}

	deleteLimit := 100
	chunks := chunkSlice(keys, deleteLimit)
	for _, chunk := range chunks {
		batchDelete(q.bm, q.bucket, chunk, &result)
	}
	return nil
}

func listAllKeys(bm *storage.BucketManager, bucket, prefix string) ([]string, error) {
	var keys []string
	marker := ""

	for {
		entries, _, nextMarker, hasNext, err := bm.ListFiles(bucket, prefix, "", marker, 1000)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			keys = append(keys, entry.Key)
		}

		if !hasNext {
			break
		}
		marker = nextMarker
	}
	return keys, nil
}

func batchDelete(bm *storage.BucketManager, bucket string, keys []string, result *DeleteResult) {
	ops := make([]string, len(keys))
	for i, key := range keys {
		ops[i] = storage.URIDelete(bucket, key)
	}

	for attempt := 1; attempt <= 3; attempt++ {
		ret, err := bm.Batch(ops)
		if err == nil {
			result.mu.Lock()
			for i, item := range ret {
				if item.Code == 200 {
					result.Success++
				} else {
					result.Failures[keys[i]] = fmt.Errorf("code:%d", item.Code)
				}
			}
			result.mu.Unlock()
			return
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
}

func chunkSlice(slice []string, size int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
