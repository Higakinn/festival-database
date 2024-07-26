package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Higakinn/festival-crawler/app/domain/models"
	"github.com/Higakinn/festival-crawler/app/domain/repository"
	"github.com/dstotijn/go-notion"
)

type FestivalRepository struct {
	// Context context.Context
	client *notion.Client
	dbName string
}

func NewFestivalRepository(client *notion.Client, dbName string) repository.FestivalRepository {
	return &FestivalRepository{client: client, dbName: dbName}
}

// FindUnPosted implements festival.FestivalRepository.
func (r FestivalRepository) FindUnPosted(ctx context.Context) ([]*models.Festival, error) {
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
	return r.queryDatabase(ctx, r.dbName, db_query)
}

func (r FestivalRepository) FindUnQuoted(ctx context.Context) ([]*models.Festival, error) {
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
	return r.queryDatabase(ctx, r.dbName, db_query)
}

// Save implements festival.FestivalRepository.
func (r FestivalRepository) Save(ctx context.Context, festival *models.Festival) error {
	r.updatePage(ctx, festival)
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

func (r *FestivalRepository) queryDatabase(ctx context.Context, databaseID string, query *notion.DatabaseQuery) ([]*models.Festival, error) {
	resp, err := r.client.QueryDatabase(ctx, databaseID, query)
	if err != nil {
		fmt.Printf("Err %v\n", err)
	}
	var results []*models.Festival
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
		fes := models.Festival{
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
		results = append(results, &fes)
	}
	return results, nil
}

func (r *FestivalRepository) updatePage(ctx context.Context, festival *models.Festival) error {
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
