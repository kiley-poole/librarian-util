package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Document struct {
	Id       int    `json:"id"`
	Filename string `json:"filename"`
	Filesize string `json:"filesize"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the directory containing the ISO files: ")
	dirName, _ := reader.ReadString('\n')
	dirName = strings.Replace(dirName, "\n", "", -1)
	wg := sync.WaitGroup{}

	files, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		wg.Add(1)
		go func(f fs.DirEntry) {
			contents, _ := os.Open(dirName + "/" + f.Name())
			res, err := postDocument(*contents)

			if err != nil {
				log.Print(err)
				return
			}

			fmt.Println(res)
			wg.Done()
		}(f)
	}
	wg.Wait()
}

func postDocument(documentIso os.File) (result string, err error) {
	documentIso.Close()
	res, err := http.Get("https://librarian-api.navsnow.com/health-check")
	if err != nil {
		log.Fatal(err)
	}
	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	result = string(response)

	return result, err
}
