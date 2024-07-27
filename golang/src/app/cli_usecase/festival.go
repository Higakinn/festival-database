package cli_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Higakinn/festival-crawler/app/domain/models"
	"github.com/Higakinn/festival-crawler/app/domain/repository"
	"github.com/Higakinn/festival-crawler/cmd"
)

type FestivalUseCase struct {
	repository   repository.FestivalRepository
	notification *cmd.XClient
}

func NewFestivalUseCase(repo repository.FestivalRepository, xClient *cmd.XClient) *FestivalUseCase {
	return &FestivalUseCase{repository: repo, notification: xClient}
}

// UnposetedList関数は 未投稿の祭り情報の一覧を取得する
func (fuc *FestivalUseCase) UnposetedList(ctx context.Context, dryRun bool) error {
	if dryRun {
		now := time.Now()
		// TODO: ドメインモデルをこちらで定義すべきではない。ファクトリを作成して、ファクトリに作成されるようにする
		festivals := []models.Festival{
			{
				Id:        "test",
				Name:      "festivalName",
				Region:    "Region",
				Access:    "Access",
				StartDate: now,
				EndDate:   now,
				Url:       "http://example.com",
			},
		}
		for i, festival := range festivals {
			fmt.Println(i + 1)
			fmt.Println(festival.GenPostContent())
		}
		return nil
	}

	// 通常の処理
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	year, month, day := time.Now().In(jst).Date()
	today := time.Date(year, month, day, 9, 0, 0, 0, jst)
	festivals, err := fuc.repository.FindByDate(ctx, today)
	if err != nil {
		return err
	}

	for i, festival := range festivals {
		fmt.Println(i + 1)
		fmt.Println(festival.GenPostContent())
	}
	return nil
}

func (fuc *FestivalUseCase) NofityUnposetedList(ctx context.Context, dryRun bool) error {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	year, month, day := time.Now().In(jst).Date()
	today := time.Date(year, month, day, 9, 0, 0, 0, jst)
	festivals, err := fuc.repository.FindByDate(ctx, today)

	if err != nil {
		return err
	}

	if festivals == nil {
		fmt.Println("未投稿の祭り情報はありません。")
		return nil
	}

	if dryRun {
		for i, festival := range festivals {
			fmt.Println(i + 1)
			fmt.Println(festival.GenPostContent())
		}
		return nil
	}

	// 通常の処理
	for i, festival := range festivals {
		fmt.Println(i + 1)
		content := festival.GenPostContent()
		postId, err := fuc.notification.Client.Post(ctx, content, festival.PosterUrl)
		if err != nil {
			return err
		}
		festival.PostId = postId
		fuc.repository.Save(ctx, festival)
	}
	return nil
}

// HoldTodayList関数は 本日開催の祭り情報を取得する
func (fuc *FestivalUseCase) HoldTodayList(ctx context.Context, dryRun bool) error {
	if dryRun {
		now := time.Now()
		// TODO: ドメインモデルをこちらで定義すべきではない。ファクトリを作成して、ファクトリに作成されるようにする
		festivals := []models.Festival{
			{
				Id:        "test",
				Name:      "festivalName",
				Region:    "Region",
				Access:    "Access",
				StartDate: now,
				EndDate:   now,
				Url:       "http://example.com",
			},
		}
		for i, festival := range festivals {
			fmt.Println(i + 1)
			fmt.Println(festival.GenQuoteRepostContent())
		}
		return nil
	}

	// 通常の処理
	// TODO: 本日開催のデータを取ってくる関数を定義すべき。
	// 現状だと、 本日開催でも投稿済みのフラグが立ってたら情報取得しないようになっている。
	isPost := false
	festivals, err := fuc.repository.FindByIsPost(ctx, isPost)
	if err != nil {
		return err
	}

	if festivals == nil {
		fmt.Println("本日開催の祭りはありません。")
	}

	for i, festival := range festivals {
		fmt.Println(i + 1)
		fmt.Println(festival.GenQuoteRepostContent())
	}
	return nil
}

// NofityHoldTodayList関数は
func (fuc *FestivalUseCase) NofityHoldTodayList(ctx context.Context, dryRun bool) error {
	isPost := false
	festivals, err := fuc.repository.FindByIsPost(ctx, isPost)
	if err != nil {
		return err
	}

	// 通常の処理
	if festivals == nil {
		fmt.Println("本日開催の祭りはありません。")
		return nil
	}

	if dryRun {
		for i, festival := range festivals {
			fmt.Println(i + 1)
			fmt.Println(festival.GenQuoteRepostContent())
		}
		return nil
	}

	for i, festival := range festivals {
		fmt.Println(i + 1)
		content := festival.GenQuoteRepostContent()
		repostId, err := fuc.notification.Client.Post(ctx, content, "")
		if err != nil {
			return err
		}
		festival.RepostId = repostId
		fuc.repository.Save(ctx, festival)
	}
	return nil
}
