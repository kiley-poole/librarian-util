package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the directory containing the ISO files: ")
	dirName, _ := reader.ReadString('\n')
	dirName = strings.Replace(dirName, "\n", "", -1)

	fmt.Print("Enter the document type uploading: ")
	docType, _ := reader.ReadString('\n')
	docType = strings.Replace(docType, "\n", "", -1)

	wg := sync.WaitGroup{}

	files, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		wg.Add(1)
		go func(f fs.DirEntry) {
			path := dirName + "/" + f.Name()
			uri := "http://localhost/api/files"

			request, err := createFormRequest(uri, path, docType)
			if err != nil {
				log.Print(err)
				return
			}

			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				log.Print(err)
				return
			} else {
				body := &bytes.Buffer{}
				_, err := body.ReadFrom(resp.Body)
				if err != nil {
					log.Print(err)
					return
				}
				resp.Body.Close()

				fmt.Println(body)
			}

			wg.Done()
		}(f)
	}
	wg.Wait()
}

// From: https://matt.aimonetti.net/posts/2013-07-golang-multipart-file-upload-example/
func createFormRequest(uri string, filePath string, docType string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	_, _ = io.Copy(part, file)
	_ = writer.WriteField("type", docType)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	return req, err
}
