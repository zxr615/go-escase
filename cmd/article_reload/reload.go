package main

import (
	"context"
	"es-demo/internal/model"
	"es-demo/pkg/es"
	"es-demo/pkg/faker"
	"log"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/pkg/errors"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/olivere/elastic/v7"
)

func main() {
	err := es.Reload(model.ArticleEsAlias, mapping(), func(newIndexName string) error {
		total := 500 // 导入总数
		batch := 100 // 每次导入数量

		var id uint32 = 0
		bar := progressbar.NewOptions(total)
		for i := 0; i < total/batch; i++ {
			bulk := es.Client.Bulk()
			for k := 0; k < batch; k++ {
				id++
				article := genArticle(id)
				bulk.Add(elastic.NewBulkIndexRequest().Index(newIndexName).Type("_doc").Id(strconv.Itoa(int(id))).Doc(article))
			}

			do, err := bulk.Do(context.Background())
			if err != nil {
				return err
			}

			if do.Errors {
				return errors.New("导入错误:do")
			}

			_ = bar.Add(batch)
		}

		_ = bar.Finish()
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func genArticle(id uint32) model.Article {
	return model.Article{
		Id:          id,
		CategoryId:  gofakeit.Uint8(),
		Title:       faker.Article(10),
		Content:     faker.Article(100),
		BrowsNum:    gofakeit.Uint8(),
		CollectNum:  gofakeit.Uint8(),
		UpvoteNum:   gofakeit.Uint8(),
		IsRecommend: uint8(gofakeit.RandomUint([]uint{1, 2, 3})),
		IsSolve:     uint8(gofakeit.RandomUint([]uint{1, 2, 3})),
		CreatedAt:   gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()).Format("2006-01-02 15:04:05"),
		UpdatedAt:   gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()).Format("2006-01-02 15:04:05"),
	}
}

func mapping() map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":           map[string]string{"type": "integer"},
				"category_id":  map[string]string{"type": "integer"},
				"title":        map[string]string{"type": "text"},
				"content":      map[string]string{"type": "text"},
				"brows_num":    map[string]string{"type": "integer"},
				"collect_num":  map[string]string{"type": "integer"},
				"upvote_num":   map[string]string{"type": "integer"},
				"is_recommend": map[string]string{"type": "integer"},
				"is_solve":     map[string]string{"type": "integer"},
				"created_at":   map[string]string{"type": "date", "format": "yyyy-MM-dd HH:mm:ss"},
				"updated_at":   map[string]string{"type": "date", "format": "yyyy-MM-dd HH:mm:ss"},
			},
		},
	}
}
