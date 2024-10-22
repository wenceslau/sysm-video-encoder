package domain

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Job struct {
	ID               string    `valid:"uuid"`
	OutputBucketPath string    `valid:"notnull"`
	Status           string    `valid:"notnull"`
	Video            *Video    `valid:"-"`
	VideoID          string    `valid:"-"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `valid:"-"`
	UpdateAt         time.Time `valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// prepare - when the function start with lowercase letter, it is private, I can't access it from another package
func (job *Job) prepare() {
	job.ID = uuid.NewV4().String()
	job.CreatedAt = time.Now()
	job.UpdateAt = time.Now()
}

func NewJob(output string, status string, video *Video) (*Job, error) {

	job := Job{
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
	}

	job.prepare()

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return &job, nil

}

// Validate - when the function start with uppercase letter, it is public, I can access it from another package
func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil
}
