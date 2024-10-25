package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type VideoUploadManager struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUploadManager() *VideoUploadManager {
	return &VideoUploadManager{}
}

func (vu *VideoUploadManager) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	// path/x/b/filename.mp4
	// split: path/x/b/
	// [0] = path/x/b/
	// [1] = filename.mp4

	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Bucket: %v\n", vu.OutputBucket)
	log.Printf("Uploading file: %v\n", objectPath)
	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	//wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUploadManager) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUploadManager) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
	}()

	countDoneWorker := 0
	for r := range returnChannel {
		countDoneWorker++

		if r != "" {
			doneUpload <- r
			break
		}
		if countDoneWorker == len(vu.Paths) {
			close(in)
		}
	}

	return nil
}

func (vu *VideoUploadManager) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)
		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("Error uploading file: %v. Error: %v", vu.Paths[x], err)
			returnChan <- err.Error()
		}
		returnChan <- ""
	}

	returnChan <- "Upload completed"
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
