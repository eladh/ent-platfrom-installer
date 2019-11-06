package utils

import (
	"bytes"
	"encoding/json"
	"github.com/kris-nova/logger"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func InvokeRequestWithContentType(url string, methodType string, data interface{} ,contentType string) string {
	return InvokeRequestInternal(url, "admin", "password", methodType, data, "" ,contentType)
}

func InvokeRequest(url string, methodType string, data interface{}) string {
	return InvokeRequestInternal(url, "admin", "password", methodType, data, "" ,"")
}

func InvokeRequestWithPassword(url string, username string, password string, methodType string, data interface{}) string {
	return InvokeRequestInternal(url, username, password, methodType, data, "" ,"")
}

func InvokeRequestWithToken(url string, methodType string, data interface{}, token string) string {
	return InvokeRequestInternal(url, "", "", methodType, data, token ,"")
}

func InvokeRequestInternal(url string, username string, password string, methodType string, data interface{}, token string ,contentType string) string {
	var request *http.Request
	var e error

	if reflect.TypeOf(data).Kind() == reflect.String {
		if FileExists(data.(string)) {
			request, e = loadFileAsset(data, methodType, url)
		} else {
			request, e = loadStringAsset(data, methodType, url)
		}
	} else {
		request, e = LoadJsonAsset(data, methodType, url ,contentType)
	}

	if e != nil {
		logger.Critical("rest error", e)
	}

	if token != "" {
		request.Header.Set("Authorization", token)
	} else {
		request.SetBasicAuth(username, password)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Critical("rest error", err)
	}

	bodyB, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//todo - only on debug mode
	logger.Always("request url is " + url + " ,data is" + string(bodyB))

	return string(bodyB)
}

func loadStringAsset(data interface{}, methodType string, url string) (*http.Request, error) {
	request, e := http.NewRequest(methodType, os.ExpandEnv(url), strings.NewReader(data.(string)))
	if e == nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return request, e
}

func LoadJsonAsset(data interface{}, methodType string, url string ,contentType string) (*http.Request, error) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}

	if contentType == "" {
		contentType = "application/json;charset=UTF-8"
	}

	request, err := http.NewRequest(methodType, os.ExpandEnv(url), bytes.NewReader(payloadBytes))
	if err == nil {
		request.Header.Set("Content-Type", contentType)
	}

	return request, err
}

func loadFileAsset(data interface{}, methodType string, url string) (*http.Request, error) {
	f, err := os.Open(data.(string))
	if err != nil {
		log.Panic(err)
	}

	return http.NewRequest(methodType, os.ExpandEnv(url), f)
}

func DownloadFile(filepath string, url string) ([]byte ,error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return nil,err
	}
	defer out.Close()

	var buf strings.Builder
	_, err = io.Copy(&buf, resp.Body)

	return []byte(buf.String()),err
}

func UploadFile(uri string, params map[string]string, paramName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return  err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("admin", "password")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Critical("rest error", err)
	}

	defer resp.Body.Close()

	return nil
}
