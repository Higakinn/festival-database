package models

import (
	"errors"
	"fmt"
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

// TODO: ドメインモデルが持つべき振る舞いではないので、適切な場所に移す
func (f *Festival) GenPostContent() string {
	date := f.StartDate.Format(time.DateOnly) + " ~ " + f.EndDate.Format(time.DateOnly)
	if f.StartDate == f.EndDate {
		date = f.StartDate.Format(time.DateOnly)
	}
	return fmt.Sprintf(`【🏮祭り情報🏮】
#%s

■ 開催期間
・%s

■ 開催場所
・%s

■ アクセス
・%s
■ 参考
%s
`,
		f.Name, date, f.Region, f.Access, f.Url)
}

// TODO: ドメインモデルが持つべき振る舞いではないので、適切な場所に移す
func (f *Festival) GenQuoteRepostContent() string {
	date := f.StartDate.Format(time.DateOnly) + " ~ " + f.EndDate.Format(time.DateOnly)
	if f.StartDate == f.EndDate {
		date = f.StartDate.Format(time.DateOnly)
	}
	return fmt.Sprintf(`【#%s】
#%s 始まります！

■ 開催期間
・%s

■ アクセス
・%s

%s
`,
		f.Region, f.Name, date, f.Access, f.XUrl)
}
