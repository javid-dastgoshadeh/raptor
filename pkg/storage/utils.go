package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
)

// NewFileUploadRequest Creates a new file upload http request with optional extra params
func NewFileUploadRequest(uri string, params map[string]string, paramName string, f *multipart.FileHeader) (*http.Request, error) {
	file, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(f.Filename))
	if err != nil {
		return nil, err
	}

	// copy
	_, err = io.Copy(part, file)

	// set params
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

// NewFileDownloadRequest Creates a new single file download http request
func NewFileDownloadRequest(uri string, id int) (*http.Request, error) {

	url := fmt.Sprintf("%s/%d", uri, id)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	return req, err
}

// NewMultipleDownloadRequest Creates a new multiple file download http request
func NewMultipleDownloadRequest(uri string, ids []int) (*http.Request, error) {
	// url?ids=1,3,5,4,7
	completeUrl := fmt.Sprintf("%s?ids=", uri)
	for _, id := range ids {
		item := fmt.Sprintf("%s,", strconv.Itoa(int(id)))
		completeUrl = completeUrl + item
	}
	req, err := http.NewRequest("GET", completeUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	return req, err
}

// NewDeleteRequest Creates a new delete file http request
func NewDeleteRequest(uri string, id int) (*http.Request, error) {
	url := fmt.Sprintf("%s/%d", uri, id)
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	return req, err
}
