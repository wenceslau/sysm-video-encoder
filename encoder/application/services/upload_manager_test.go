package services_test

import (
	"encoder/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("sysm_catalogo_videos/project-299e23d88f8643bd9acf1833de39144b")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUploadManager()
	videoUpload.OutputBucket = "sysm_catalogo_videos"
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(100, doneUpload)

	result := <-doneUpload
	require.Equal(t, "Upload completed", result)

	err = videoService.Finish()
	require.Nil(t, err)

}
