package main

import (
	"es-demo/internal/api"
	"es-demo/pkg/es"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	exec(func() error {
		// 初始化 es
		if err := es.New(); err != nil {
			return err
		}

		// 初始化路由
		router := gin.Default()
		v1 := router.Group("/article/v1")
		{
			v1.GET("/search", api.NewArticleV1().Search)
			v1.GET("/recommend", api.NewArticleV1().Recommend)
			v1.GET("/related", api.NewArticleV1().Related)
		}

		v2 := router.Group("/article/v2")
		{
			v2.GET("/search", api.NewArticleV2().Search)
			v2.GET("/recommend", api.NewArticleV2().Recommend)
			v2.GET("/related", api.NewArticleV2().Related)
		}

		if err := router.Run(); err != nil {
			return err
		}

		return nil
	})
}

func exec(f func() error) {
	err := f()
	if err != nil {
		log.Fatal(err)
	}
}
