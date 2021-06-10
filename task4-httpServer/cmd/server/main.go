package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

var (
	fileSystem = "/tmp/upload"
)

type File struct {
	Name string `json:"filename"`
	Size int    `json:"sizeByte"`
}

func main() {
	uploadHandler := &UploadHandler{
		UploadDir: fileSystem,
	}
	http.Handle("/upload", uploadHandler)
	http.Handle("/ls", uploadHandler)
	http.Handle("/", uploadHandler)

	go http.ListenAndServe(":8000", nil)

	dirToServe := http.Dir(uploadHandler.UploadDir)

	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fs.ListenAndServe()
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ext string
	fileList, err := PrintFileSystem(os.DirFS(fileSystem))
	if r.URL.Path == "/" {
		var selectedFiles []File
		param := r.FormValue("ext")
		if param == "txt" {
			ext = "txt"
		}
		if param == "jpg" {
			ext = "jpg"
		}
		for _, item := range fileList {
			fileName := strings.Split(item.Name, ".")
			if fileName[1] == ext {
				selectedFiles = append(selectedFiles, item)
			}
		}
		jsonOut, err := json.Marshal(selectedFiles)
		if err != nil {
			http.Error(w, "Unable to do marshall to Json", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, string(jsonOut))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err != nil {
			http.Error(w, "Unable to get content of root directory: "+fileSystem, http.StatusBadRequest)
			return
		}
		jsonOut, err := json.Marshal(fileList)
		if err != nil {
			http.Error(w, "Unable to do marshall to Json", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, string(jsonOut))
		return
	case http.MethodPost:
		h.HostAddr = "http://localhost:8080"
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}

		timeStamp := strconv.Itoa(int(time.Now().UnixNano() / 1000000))
		filePath := h.UploadDir + "/" + timeStamp + "-" + header.Filename

		err = ioutil.WriteFile(filePath, data, 0777)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		fileLink := h.HostAddr + "/" + timeStamp + "-" + header.Filename
		fmt.Fprintln(w, fileLink)
	}
}

func PrintFileSystem(filepath fs.FS) ([]File, error) {
	var fileList []File
	var file File
	err := fs.WalkDir(filepath, ".", func(path string, info fs.DirEntry, err error) error {
		fileInfo, _ := info.Info()
		if !info.IsDir() {
			file.Name = fileInfo.Name()
			file.Size = int(fileInfo.Size())
			fileList = append(fileList, file)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error occurs: %v", err)
	}
	return fileList, nil
}
