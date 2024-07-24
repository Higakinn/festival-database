package festival

import (
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
	XUrl      string
}

func (f Festival) PostContent() string {
	date := f.StartDate.Format(time.DateOnly) + " ~ " + f.EndDate.Format(time.DateOnly)
	if f.StartDate == f.EndDate {
		date = f.StartDate.Format(time.DateOnly)
	}
	fmt.Println(f.PosterUrl)
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

%s
`,
		f.Name, date, f.Region, f.Access, f.Url, f.PosterUrl)
}
