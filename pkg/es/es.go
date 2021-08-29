package es

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

var Client *elastic.Client

func New() error {
	client, err := elastic.NewSimpleClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		return err
	}

	Client = client

	return nil
}

// CreateIndex 创建一个新索引
func CreateIndex(aliasIndexName string, mapping map[string]interface{}) (*elastic.IndicesCreateResult, error) {
	// 索引加上当前时间,eg:article_20210829185320
	timeNewIndexName := aliasIndexName + "_" + time.Now().Format("20060102150405")

	return Client.
		CreateIndex(timeNewIndexName).
		BodyJson(mapping).
		Do(context.Background())
}

// SetAlias 为索引设置别名
func SetAlias(indexName string, aliasName string) (*elastic.AliasResult, error) {
	ctx := context.Background()

	// 获取别名下的所有索引
	catAliases, err := Client.CatAliases().Alias(aliasName).Do(ctx)
	if err != nil {
		return nil, err
	}

	// 新增新索引 & 删除旧索引
	// 因为一个别名可以关联多个 index，又因为设计的是一个索引关联一个别名，所以循环增加删除也就不存在什么问题了。
	alias := Client.Alias()
	alias.Action(elastic.NewAliasAddAction(aliasName).Index(indexName))
	for _, as := range catAliases {
		alias.Action(elastic.NewAliasRemoveAction(aliasName).Index(as.Index))
	}

	return alias.Do(ctx)
}

// DeleteIndex 删除索引
func DeleteIndex(indices ...string) (*elastic.IndicesDeleteResponse, error) {
	return Client.DeleteIndex(indices...).Do(context.Background())
}
