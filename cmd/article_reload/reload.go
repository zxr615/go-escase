package main

import (
	"context"
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
	type Article struct {
		Id          uint32 `json:"id"`           // Id
		CategoryId  uint8  `json:"category_id"`  // 分类
		Title       string `json:"title"`        // 标题
		Content     string `json:"content"`      // 内容
		BrowsNum    uint8  `json:"brows_num"`    // 浏览量
		CollectNum  uint8  `json:"collect_num"`  // 收藏量
		UpvoteNum   uint8  `json:"upvote_num"`   // 点赞量
		IsRecommend uint8  `json:"is_recommend"` // 是否推荐:1=是;2=否
		IsSolve     uint8  `json:"is_solve"`     // 是否解决:1=是;2=否
		CreatedAt   string `json:"created_at"`   // 创建时间
		UpdatedAt   string `json:"updated_at"`   // 更新时间
	}

	err := es.Reload("article", mapping(), func(newIndexName string) error {
		total := 500
		batch := 100

		var id uint32 = 0
		bar := progressbar.NewOptions(total)
		for i := 0; i < total/batch; i++ {
			bulk := es.Client.Bulk()
			for k := 0; k < batch; k++ {
				id++
				article := Article{
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

	log.Println("SUCCESS")
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
