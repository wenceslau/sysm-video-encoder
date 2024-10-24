package services

import (
	"cloud.google.com/go/storage"
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {

	// Context is used to handle the request deadline, I can cancel the request after a certain time if I want
	ctx := context.Background()

	// Create a new client to interact with the Google Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// Create a new bucket using the client, and get the file from the bucket
	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)

	// Create a new reader to read the file from the bucket
	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	// Read the content of the file
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// Create a new file to save the content read from the bucket
	f, err := os.Create(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	// Write the content read from the bucket to the file
	_, err = f.Write(body)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Video %v has been downloaded", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	// Create a new folder to save the fragments of the video. os.ModePerm is the permission to access the folder
	err := os.Mkdir(os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	// Path to the mp4dash binary
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID+".frag")

	// Command to segment the video
	cmdArgs = append(cmdArgs, "--use-segment-timeline")

	// Command
	cmdArgs = append(cmdArgs, "-o")

	// Path to save the video
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID)

	// Command
	cmdArgs = append(cmdArgs, "-f")

	// Path where the command will be executed
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	// Execute the command with all the arguments
	cmd := exec.Command("mp4dash", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Println("Error removing the video: ", err)
		return err
	}

	err = os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag")
	if err != nil {
		log.Println("Error removing the frag: ", err)
		return err
	}

	err = os.RemoveAll(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID)
	if err != nil {
		log.Println("Error removing the folder with the fragments: ", err)
		return err
	}

	log.Println("Files have been removed", v.Video.ID)

	return err

}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("Output: %s\n", output)
	}
}

// bucketName / Folder: project-299e23d88f8643bd9acf1833de39144b/
// Path / File name:	resource-7454ac24dfaf4f639c0c645c0e5b5086-file_example_mp4_640_3mg.mp4
// Path / File name:	resource-5fe9b506258d4b8290178ee01f06df38-file_example_mp4_1920_18mg.mp4
