package domain_test

import (
	// 'uuid' os the alias for the github.com/satori/go.uuid package
	uuid "github.com/satori/go.uuid"

	"encoder/domain"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidateIfVideIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()
	require.Error(t, err)
}

func TestVideIdIsNotUUID(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "abc"
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Error(t, err)
}

func TestVideValidation(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Nil(t, err)
}
