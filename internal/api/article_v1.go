package api

import (
	"context"
	"encoding/json"
	"es-demo/internal/model"
	"es-demo/pkg/es"
	"log"
	"net/http"

	"github.com/olivere/elastic/v7"

	"github.com/gin-gonic/gin"
)

type ArticleV1 struct{}

func NewArticleV1() *ArticleV1 {
	return &ArticleV1{}
}

const (
	SortBrowseDesc  = iota + 1 // 浏览量倒序
	SortCollectDesc            // 收藏倒序
	SortUpvoteDesc             // 点赞倒序
)

// Search 文章搜索
func (a ArticleV1) Search(c *gin.Context) {
	req := new(model.SearchRequest)
	if err := c.ShouldBind(req); err != nil {
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

	// 是否解决
	if req.IsSolve != 0 {
		builder.Query(bq.Filter(elastic.NewTermQuery("is_solve", req.IsSolve)))
	}

	// 排序
	switch req.Sort {
	case SortBrowseDesc:
		builder.Sort("brows_num", false)
	case SortUpvoteDesc:
		builder.Sort("upvote_num", false)
	case SortCollectDesc:
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

	// 序列化数据
	list := make([]model.SearchResponse, len(do.Hits.Hits))
	for i, raw := range do.Hits.Hits {
		tmpArticle := model.SearchResponse{}
		if err := json.Unmarshal(raw.Source, &tmpArticle); err != nil {
			log.Println(err)
		}

		list[i] = tmpArticle
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  list,
		},
	})
}

// Recommend 文章推荐
func (a ArticleV1) Recommend(c *gin.Context) {
	// 构建搜索
	builder := es.Client.Search().Index(model.ArticleEsAlias)

	bq := elastic.NewBoolQuery()

	builder.Query(bq.Filter(
		// 推荐文章
		elastic.NewTermQuery("category_id", model.ArticleIsRecommendYes),
		// 已解决
		elastic.NewTermQuery("is_solve", model.ArticleIsSolveYes),
	))

	// 浏览量排序
	builder.Sort("brows_num", false)

	do, err := builder.From(0).Size(10).Do(context.Background())
	if err != nil {
		return
	}

	// 序列化数据
	list := make([]model.RecommendResponse, len(do.Hits.Hits))
	for i, raw := range do.Hits.Hits {
		tmpArticle := model.RecommendResponse{}
		if err := json.Unmarshal(raw.Source, &tmpArticle); err != nil {
			log.Println(err)
		}

		list[i] = tmpArticle
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": len(list),
			"list":  list,
		},
	})
}

// Related 相关文章
func (a ArticleV1) Related(c *gin.Context) {

}
