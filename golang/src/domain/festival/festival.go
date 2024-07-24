package festival

import (
	"errors"
	"time"
)

type Festival struct {
	Id        string
	Name      string
	Region    string
	Access    string
	StartDate time.Time
	EndDate   time.Time
	Url       string
	PosterUrl string
	PostId    string
	RepostId  string
	XUrl      string
}

func (f *Festival) Validate() error {
	if f.Id == "" {
		return errors.New("id not empty")
	}
	if f.Name == "" {
		return errors.New("name not empty")
	}
	if f.Region == "" {
		return errors.New("region not empty")
	}
	if f.Url == "" {
		return errors.New("	Url not empty")
	}
	return nil
}
