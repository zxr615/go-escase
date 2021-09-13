package api

import (
	"es-demo/internal/model"
	"es-demo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleV2 struct{}

func NewArticleV2() *ArticleV2 {
	return &ArticleV2{}
}

// Search 文章搜索
func (a ArticleV2) Search(c *gin.Context) {
	req := new(model.SearchRequest)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(400, err.Error())
		return
	}

	list, total, err := service.NewArticle().
		WhereKeyword(req.Keyword).
		WhereCategoryId(req.CategoryId).
		WhereIsSolve(req.IsSolve).
		Sort(req.Sort).
		Paginate(req.Page, req.PageSize).
		DecodeSearch()

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  list,
		},
	})
	return
}

// Recommend 文章推荐
func (a ArticleV2) Recommend(c *gin.Context) {
	list, _, err := service.NewArticle().
		WhereIsRecommend(model.ArticleIsRecommendYes).
		WhereIsSolve(model.ArticleIsSolveYes).
		OrderByDesc("brows_num").
		PageSize(10).
		DecodeRecommend()

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": len(list),
			"list":  list,
		},
	})
	return
}

// Related 相关文章
func (a ArticleV2) Related(c *gin.Context) {
	req := new(model.RelatedRequest)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(400, err.Error())
		return
	}

	list, _, err := service.NewArticle().
		WhereCategoryId(req.CategoryId).
		WhereIsRecommend(model.ArticleIsRecommendYes).
		WhereIsSolve(model.ArticleIsSolveYes).
		OrderByDesc("brows_num").
		PageSize(10).
		DecodeRelated()
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list": list,
		},
	})
}
