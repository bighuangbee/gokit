package function

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var RequestTimeout = time.Second * 3

func HttpPost(addr string, params map[string]interface{}) ([]byte, error) {
	jsonData, _ := json.Marshal(&params)
	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: RequestTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type FormFile struct {
	FormField string
	Filename  string
}

func HttpFormPost(addr string, params map[string]string, formfile *FormFile) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// formFile
	fileReader, err := os.Open(formfile.Filename)
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	part, err := writer.CreateFormFile(formfile.FormField, formfile.Filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, fileReader)
	if err != nil {
		return nil, err
	}

	for key, value := range params {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", addr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: RequestTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
