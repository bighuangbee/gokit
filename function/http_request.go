package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var RequestTimeout = time.Second * 3

func HttpRequest(method string, url string, header map[string]interface{}, body map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化 JSON 失败: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for key, value := range header {
		req.Header.Add(key, fmt.Sprintf("%v", value))
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
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
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
