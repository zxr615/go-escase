package main

import (
	"context"
	"es-demo/internal/model"
	"es-demo/pkg/es"
	"es-demo/pkg/faker"
	"log"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/schollz/progressbar/v3"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	// 初始化 es
	if err := es.New(); err != nil {
		log.Fatalln(err)
	}

	//type Article struct {
	//	Id          uint32 `json:"id"`           // Id
	//	CategoryId  uint8  `json:"category_id"`  // 分类
	//	Title       string `json:"title"`        // 标题
	//	Content     string `json:"content"`      // 内容
	//	BrowsNum    uint8  `json:"brows_num"`    // 浏览量
	//	CollectNum  uint8  `json:"collect_num"`  // 收藏量
	//	UpvoteNum   uint8  `json:"upvote_num"`   // 点赞量
	//	IsRecommend uint8  `json:"is_recommend"` // 是否推荐:1=是;2=否
	//	IsSolve     uint8  `json:"is_solve"`     // 是否解决:1=是;2=否
	//	CreatedAt   string `json:"created_at"`   // 创建时间
	//	UpdatedAt   string `json:"updated_at"`   // 更新时间
	//}

	total := 500
	batch := 100
	createIndexRs, err := es.CreateIndex("article", mapping())
	if err != nil {
		log.Fatalln(err)
	}

	var id uint32 = 0
	bar := progressbar.NewOptions(total)
	for i := 0; i < total/batch; i++ {
		bulk := es.Client.Bulk()
		for k := 0; k < batch; k++ {
			id++
			article := genArticle(id)
			bulk.Add(elastic.NewBulkIndexRequest().Index(createIndexRs.Index).Type("_doc").Id(strconv.Itoa(int(id))).Doc(article))
		}

		do, err := bulk.Do(context.Background())
		if err != nil {
			log.Fatalln(err, do)
		}

		if do.Errors {
			log.Fatalln("导入错误")
		}

		time.Sleep(5 * time.Millisecond)
		_ = bar.Add(batch)
	}

	alias, err := es.SetAlias(createIndexRs.Index, "article")
	if err != nil {
		log.Fatalln(err)
	}

	if !alias.Acknowledged {
		log.Fatalln("导入错误")
	}

	_ = bar.Finish()
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
