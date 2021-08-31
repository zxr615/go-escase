package api

import (
	"context"
	"encoding/json"
	"es-demo/internal/model"
	"es-demo/pkg/es"
	"log"

	"github.com/olivere/elastic/v7"

	"github.com/gin-gonic/gin"
)

type ArticleV1 struct{}

func NewArticleV1() *ArticleV1 {
	return &ArticleV1{}
}

type SearchRequest struct {
	Keyword    string `form:"keyword"`                              // 关键词
	CategoryId uint8  `form:"category_id"`                          // 分类
	Sort       uint8  `form:"sort" binding:"omitempty,oneof=1 2 3"` // 排序 1=浏览量；2=收藏；3=点赞；
	isSolve    uint8  `form:"is_solve"`                             // 是否解决
	Page       int    `form:"page,default=1"`
	PageSize   int    `json:"page_size,default=10"`
}

type CommonDataResponse struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
type CommonResponse struct {
	Code int                `json:"code"`
	Data CommonDataResponse `json:"data"`
}

const (
	SortBrowse = iota + 1
	SortCollect
	SortUpvote
)

// Search 文章搜索
func (a ArticleV1) Search(c *gin.Context) {
	req := new(SearchRequest)
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(500, err)
	}

	// 构建搜索
	builder := es.Client.Search().Index(model.ArticleEsAlias)

	bq := elastic.NewBoolQuery()
	// 标题
	if req.Keyword != "" {
		builder.Query(bq.Must(elastic.NewMatchQuery("title", req.Keyword)))
	}

	// 分类
	if req.CategoryId != 0 {
		builder.Query(bq.Filter(elastic.NewTermQuery("category_id", req.CategoryId)))
	}

	// 排序
	switch req.Sort {
	case SortBrowse:
		builder.Sort("brows_num", false)
	case SortUpvote:
		builder.Sort("upvote_num", false)
	case SortCollect:
		builder.Sort("collect_num", false)
	default:
		builder.Sort("created_at", false)
	}

	// 分页
	from := (req.Page - 1) * req.PageSize
	builder.
		FetchSourceContext(elastic.NewFetchSourceContext(true)).
		From(from).
		Size(req.PageSize)

	// 执行查询
	do, err := builder.Do(context.Background())
	if err != nil {
		c.JSON(500, err)
	}

	// 获取匹配到的数量
	total := do.TotalHits()

	list := make([]model.Article, len(do.Hits.Hits))
	for i, raw := range do.Hits.Hits {
		tmpArticle := new(model.Article)
		if err := json.Unmarshal(raw.Source, tmpArticle); err != nil {
			log.Println(err)
		}

		list[i] = *tmpArticle
	}

	c.JSON(200, CommonResponse{
		Code: 200,
		Data: CommonDataResponse{
			Total: total,
			List:  list,
		},
	})
}

// Recommend 文章推荐
func (a ArticleV1) Recommend(c *gin.Context) {

}

// Related 相关文章
func (a ArticleV1) Related(c *gin.Context) {

}
