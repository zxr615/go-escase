package service

import (
	"context"
	"encoding/json"
	"es-demo/internal/model"
	"es-demo/pkg/es"
	"es-demo/pkg/structer"
	"log"

	"github.com/olivere/elastic/v7"
)

type article struct {
	must   []elastic.Query
	filter []elastic.Query
	sort   []elastic.Sorter
	from   int
	size   int
}

const (
	SortBrowseDesc  = iota + 1 // 浏览量倒序
	SortCollectDesc            // 收藏倒序
	SortUpvoteDesc             // 点赞倒序
)

func NewArticle() *article {
	return &article{
		must:   make([]elastic.Query, 0),
		filter: make([]elastic.Query, 0),
		sort:   make([]elastic.Sorter, 0),
		from:   0,
		size:   10,
	}
}

// WhereKeyword 关键词
func (a article) WhereKeyword(keyword string) article {
	if keyword != "" {
		a.must = append(a.must, elastic.NewMatchQuery("title", keyword))
	}

	return a
}

// WhereCategoryId 分类
func (a article) WhereCategoryId(categoryId uint8) article {
	if categoryId != 0 {
		a.filter = append(a.filter, elastic.NewTermQuery("category_id", categoryId))
	}

	return a
}

// WhereIsSolve 是否已解决
func (a article) WhereIsSolve(isSolve uint8) article {
	if isSolve != 0 {
		a.filter = append(a.filter, elastic.NewTermQuery("is_solve", isSolve))
	}

	return a
}

// WhereIsRecommend 是否推荐
func (a article) WhereIsRecommend(isRecommend uint8) article {
	if isRecommend != 0 {
		a.filter = append(a.filter, elastic.NewTermQuery("is_recommend", isRecommend))
	}

	return a
}

// Sort 排序
func (a article) Sort(sort uint8) article {
	switch sort {
	case SortBrowseDesc:
		return a.OrderByDesc("brows_num")
	case SortUpvoteDesc:
		return a.OrderByDesc("upvote_num")
	case SortCollectDesc:
		return a.OrderByDesc("collect_num")
	}

	return a
}

// OrderByDesc 通过字段倒序排序
func (a article) OrderByDesc(field string) article {
	a.sort = append(a.sort, elastic.SortInfo{Field: field, Ascending: false})
	return a
}

// OrderByAsc 通过字段正序排序
func (a article) OrderByAsc(field string) article {
	a.sort = append(a.sort, elastic.SortInfo{Field: field, Ascending: true})
	return a
}

// Paginate 分页
// page 当前页码
// pageSize 每页数量
func (a article) Paginate(page, pageSize int) article {
	a.from = (page - 1) * pageSize
	a.size = pageSize
	return a
}

// PageSize 每页数量
func (a article) PageSize(pageSize int) article {
	a.size = pageSize
	return a
}

// DecodeSearch 搜索结果
func (a article) DecodeSearch() ([]model.SearchResponse, int64, error) {
	rawList, total, err := a.Searcher(new(model.SearchResponse))
	if err != nil {
		return nil, total, err
	}

	list := make([]model.SearchResponse, len(rawList))
	for i, raw := range rawList {
		tmp := model.SearchResponse{}
		if err := json.Unmarshal(raw, &tmp); err != nil {
			log.Println(err)
			continue
		}

		list[i] = tmp
	}

	return list, total, nil
}

// DecodeRecommend 推荐结果
func (a article) DecodeRecommend() ([]model.RecommendResponse, int64, error) {
	rawList, total, err := a.Searcher(new(model.RecommendResponse))
	if err != nil {
		return nil, total, err
	}

	list := make([]model.RecommendResponse, len(rawList))
	for i, raw := range rawList {
		tmp := model.RecommendResponse{}
		if err := json.Unmarshal(raw, &tmp); err != nil {
			log.Println(err)
			continue
		}

		list[i] = tmp
	}

	return list, total, nil
}

// DecodeRelated decode 相关结果
func (a article) DecodeRelated() ([]model.RelatedResponse, int64, error) {
	rawList, total, err := a.Searcher(new(model.RelatedResponse))
	if err != nil {
		return nil, total, err
	}

	list := make([]model.RelatedResponse, len(rawList))
	for i, raw := range rawList {
		tmp := model.RelatedResponse{}
		if err := json.Unmarshal(raw, &tmp); err != nil {
			log.Println(err)
			continue
		}

		list[i] = tmp
	}

	return list, total, nil
}

// Searcher 执行查询
func (a article) Searcher(include ...interface{}) ([]json.RawMessage, int64, error) {
	builder := es.Client.Search().Index(model.ArticleEsAlias)

	// 查询的字段
	includeKeys := make([]string, 0)
	if len(include) > 0 {
		includeKeys = structer.Keys(include[0])
	}

	// 构建查询
	builder.Query(
		// 构建 bool query 条件
		elastic.NewBoolQuery().Must(a.must...).Filter(a.filter...),
	)

	// 执行查询
	do, err := builder.
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include(includeKeys...)).
		From(a.from).
		Size(a.size).
		SortBy(a.sort...).
		Do(context.Background())

	if err != nil {
		return nil, 0, err
	}

	total := do.TotalHits()
	list := make([]json.RawMessage, len(do.Hits.Hits))
	for i, hit := range do.Hits.Hits {
		list[i] = hit.Source
	}

	return list, total, nil
}
