package domain

import (
	"github.com/asaskevich/govalidator"
	"time"
)

type Video struct {
	ID         string    `valid:"uuid"`
	ResourceID string    `valid:"notnull"`
	FilePath   string    `valid:"notnull"`
	CreatedAt  time.Time `valid:"-"`
}

func init() {
	// SetFieldsRequiredByDefault - the function sets the default behavior of the govalidator package to require all fields to be present
	govalidator.SetFieldsRequiredByDefault(true)
}

// NewVideo - the '*' indicates that the function returns a pointer to a Video struct
func NewVideo() *Video {
	// the '&' indicates that the variable is a pointer
	return &Video{}
}

// Validate - the '(video *Video)' indicates that the function is a method of the Video struct
func (video *Video) Validate() error {
	// the '_' indicates that the variable is not being used
	_, err := govalidator.ValidateStruct(video)

	if err != nil {
		return err
	}

	return nil
}
