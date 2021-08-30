package api

import "github.com/gin-gonic/gin"

type ArticleV1 struct{}

func NewArticleV1() *ArticleV1 {
	return &ArticleV1{}
}

// Search 文章搜索
func (a ArticleV1) Search(c *gin.Context) {

}

// Recommend 文章推荐
func (a ArticleV1) Recommend(c *gin.Context) {

}

// Related 相关文章
func (a ArticleV1) Related(c *gin.Context) {

}
