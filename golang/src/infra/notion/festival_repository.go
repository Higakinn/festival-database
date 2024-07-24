package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	festival "github.com/Higakinn/festival-crawler/domain/festival"
	"github.com/dstotijn/go-notion"
	gn "github.com/dstotijn/go-notion"
)

type FestivalRepository struct {
	Context context.Context
	client  *gn.Client
	dbName  string
}

// FindUnPosted implements festival.FestivalRepository.
func (r FestivalRepository) FindUnPosted() ([]festival.Festival, error) {
	db_query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			And: []notion.DatabaseQueryFilter{
				{
					Property: "is_post",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Checkbox: &notion.CheckboxDatabaseQueryFilter{
							Equals: &[]bool{false}[0],
						},
					},
				},
			},
		},
	}
	return r.queryDatabase(r.Context, r.dbName, db_query)
}

func (r FestivalRepository) FindUnQuoted() ([]festival.Festival, error) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	year, month, day := time.Now().In(jst).Date()
	today := time.Date(year, month, day, 9, 0, 0, 0, jst)
	fmt.Println(today)

	db_query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			And: []notion.DatabaseQueryFilter{
				{
					Property: "is_post",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Checkbox: &notion.CheckboxDatabaseQueryFilter{
							Equals: &[]bool{true}[0],
						},
					},
				},
				{
					Property: "is_repost",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Checkbox: &notion.CheckboxDatabaseQueryFilter{
							Equals: &[]bool{false}[0],
						},
					},
				},
				{
					Property: "date",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Date: &notion.DatePropertyFilter{
							Equals: &today,
						},
					},
				},
			},
		},
	}
	return r.queryDatabase(r.Context, r.dbName, db_query)
}

// Save implements festival.FestivalRepository.
func (r FestivalRepository) Save(festival festival.Festival) error {
	r.updatePage(r.Context, &festival)
	return nil
}

type notionDBProperty struct {
	Access       notion.DatabasePageProperty `json:"access"`
	Date         notion.DatabasePageProperty `json:"date"`
	FestivalName notion.DatabasePageProperty `json:"festival_name"`
	IsPost       notion.DatabasePageProperty `json:"is_post"`
	Link         notion.DatabasePageProperty `json:"link"`
	PostID       notion.DatabasePageProperty `json:"post_id,omitempty"`
	ReostID      notion.DatabasePageProperty `json:"repost_id,omitempty"`
	Region       notion.DatabasePageProperty `json:"region"`
	XUrl         notion.DatabasePageProperty `json:"x url,omitempty"`
	Poster       notion.DatabasePageProperty `json:"poster"`
	StrCount     notion.DatabasePageProperty `json:"str_count,omitempty"`
}

func NewFestivalRepository(ctx context.Context, apiKey string, dbName string) FestivalRepository {
	return FestivalRepository{Context: ctx, client: gn.NewClient(apiKey), dbName: dbName}
}

func (r *FestivalRepository) queryDatabase(ctx context.Context, databaseID string, query *gn.DatabaseQuery) ([]festival.Festival, error) {
	resp, err := r.client.QueryDatabase(ctx, databaseID, query)
	if err != nil {
		fmt.Printf("Err %v\n", err)
	}
	var results []festival.Festival
	for _, page := range resp.Results {
		jsonBytes, err := json.Marshal(page.Properties)
		if err != nil {
			// エラー処理
			fmt.Println("JSON marshal error:", err)
		}

		var n notionDBProperty
		if err := json.Unmarshal(jsonBytes, &n); err != nil {
			fmt.Println(err)
		}

		endDate := n.Date.Date.Start.Time
		if n.Date.Date.End != nil {
			endDate = n.Date.Date.End.Time
		}

		posterUrl := ""
		if len(n.Poster.Files) > 0 {
			posterUrl = n.Poster.Files[0].External.URL
		}

		postId := ""
		if len(n.PostID.RichText) > 0 {
			postId = n.PostID.RichText[0].PlainText
		}

		repostId := ""
		if len(n.ReostID.RichText) > 0 {
			repostId = n.ReostID.RichText[0].PlainText
		}

		xUrl := ""
		if n.XUrl.Formula != nil {
			xUrl = *n.XUrl.Formula.String
		}
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		fes := festival.Festival{
			Id:        page.ID,
			Name:      n.FestivalName.Title[0].PlainText,
			Region:    n.Region.RichText[0].PlainText,
			Access:    n.Access.RichText[0].PlainText,
			StartDate: n.Date.Date.Start.Time.In(jst),
			EndDate:   endDate.In(jst),
			Url:       *n.Link.URL,
			PosterUrl: posterUrl,
			PostId:    postId,
			RepostId:  repostId,
			XUrl:      xUrl,
		}
		if err := fes.Validate(); err != nil {
			log.Fatal(err)
			return nil, err
		}
		results = append(results, fes)
	}
	return results, nil
}

func (r *FestivalRepository) updatePage(ctx context.Context, festival *festival.Festival) error {
	isPost := festival.PostId != ""
	isRepost := festival.RepostId != ""
	params := &notion.UpdatePageParams{
		DatabasePageProperties: notion.DatabasePageProperties{
			"post_id": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{Text: &notion.Text{Content: festival.PostId}},
				},
			},
			"is_post": notion.DatabasePageProperty{
				Checkbox: &[]bool{isPost}[0],
			},
			"repost_id": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{Text: &notion.Text{Content: festival.RepostId}},
				},
			},
			"is_repost": notion.DatabasePageProperty{
				Checkbox: &[]bool{isRepost}[0],
			},
		},
	}
	_, err := r.client.UpdatePage(ctx, festival.Id, *params)
	if err != nil {
		return err
	}
	return nil
}
