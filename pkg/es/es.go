package es

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/olivere/elastic/v7"
)

var Client *elastic.Client

func New() error {
	esLogFile, err := os.OpenFile("es.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	client, err := elastic.NewSimpleClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetInfoLog(log.New(esLogFile, "", log.Lshortfile|log.LstdFlags)),
		elastic.SetTraceLog(log.New(esLogFile, "", log.Lshortfile|log.LstdFlags)),
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
func DeleteIndex(indices []string) (*elastic.IndicesDeleteResponse, error) {
	if len(indices) == 0 {
		return nil, nil
	}
	return Client.DeleteIndex().Index(indices).Do(context.Background())
}

// DeleteIndexByPrefix 通过前缀批量删除索引
func DeleteIndexByPrefix(prefix string, excludeIndex string) (*elastic.IndicesDeleteResponse, error) {
	indices, err := Client.CatIndices().Index(prefix + "*").Do(context.Background())
	if err != nil {
		return nil, err
	}

	var wait []string
	for _, i := range indices {
		if i.Index == excludeIndex {
			continue
		}

		wait = append(wait, i.Index)
	}

	return DeleteIndex(wait)
}

// Reload 封装导入流程：1 创建新索引；2 导入数据；3 切换别名；4 删除旧索引
// aliasName 别名
// mapping es的 mapping 定义
// callback 导入数据具体实现
func Reload(aliasName string, mapping map[string]interface{}, callback func(newIndexName string) error) error {
	if err := New(); err != nil {
		log.Fatalln(err)
	}

	// 创建一个新的索引
	createIndexRs, err := CreateIndex(aliasName, mapping)
	if err != nil {
		return err
	}
	createNewIndexName := createIndexRs.Index

	// 调用自定义执行程序
	if err := callback(createNewIndexName); err != nil {
		return err
	}

	// 切换别名
	alias, err := SetAlias(createNewIndexName, aliasName)
	if err != nil {
		return err
	}
	if !alias.Acknowledged {
		return errors.New("导入错误")
	}

	// 删除旧索引
	if _, err := DeleteIndexByPrefix(aliasName, createNewIndexName); err != nil {
		return err
	}

	return nil
}
