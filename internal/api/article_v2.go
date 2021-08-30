package api

import "github.com/gin-gonic/gin"

type ArticleV2 struct{}

func NewArticleV2() *ArticleV2 {
	return &ArticleV2{}
}

// Search 文章搜索
func (a ArticleV2) Search(c *gin.Context) {

	c.JSON(200, "test")
}

// Recommend 文章推荐
func (a ArticleV2) Recommend(c *gin.Context) {

}

// Related 相关文章
func (a ArticleV2) Related(c *gin.Context) {

}
