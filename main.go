package main

import (
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
		r := gin.Default()
		if err := r.Run(); err != nil {
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
