package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewVideoRepository(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	_, err := repo.Insert(video)
	require.Nil(t, err)

	v, err := repo.Find(video.ID)
	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, video.ID, v.ID)
}
