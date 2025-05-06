package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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

func RequestForm(method, urlStr string, fields map[string]interface{}, header map[string]interface{}) ([]byte, error) {
	var req *http.Request
	var err error

	if method == "GET" {
		// 构造查询参数
		qs := url.Values{}
		for key, value := range fields {
			qs.Add(key, fmt.Sprintf("%v", value))
		}
		urlStr = urlStr + "?" + qs.Encode()
		req, err = http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return nil, err
		}
	} else if method == "POST" {
		// 创建 multipart/form-data 请求体
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for key, value := range fields {
			err := writer.WriteField(key, fmt.Sprintf("%v", value))
			if err != nil {
				return nil, err
			}
		}

		// 关闭 writer 以结束 multipart 数据
		err = writer.Close()
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest("POST", urlStr, body)
		if err != nil {
			return nil, err
		}

		// 使用 writer 自动生成的 Content-Type
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		return nil, errors.New("unsupported HTTP method")
	}

	// 添加请求头
	for key, value := range header {
		req.Header.Add(key, fmt.Sprintf("%v", value))
	}

	// 发起 HTTP 请求
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("请求出错, " + err.Error())
	}

	return respBody, nil
}
