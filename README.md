## clone 
```shell
git clone git@github.com:zxr615/go-escase.git
```

## go mod
```shell
go mod download
```

## 导入测试文章数据
> 内容是鲁迅的《孔乙己》，随机截取文本内容作为文章标题与文章内容
```shell
go run go-escase/cmd/article_reload
```

## 启动程序
```shell
go run main.go
```

## 分词
使用 [ik_smart 分词](https://github.com/medcl/elasticsearch-analysis-ik)

## api 描述
search    搜索文章：按照给定条件搜索结果  
related   相关文章：根据给定分类 id 搜索当前分类下已解决的推荐文章，浏览量从高到低排序  
recommend 推荐文章：搜索已解决的推荐文章，浏览量从高到低排序  

## api-v1
搜索、相关、推荐文章  
`curl GET '127.0.0.1:8080/article/v1/search?keyword=茴香豆&page=1&page_size=5&sort=1'`  
`curl GET '127.0.0.1:8080/article/v1/related?category_id=1'`  
`curl GET '127.0.0.1:8080/article/v1/recommend'`  

## api-v1
搜索、相关、推荐文章  
`curl GET '127.0.0.1:8080/article/v2/search?keyword=茴香豆&page=1&page_size=5&sort=1'`  
`curl GET '127.0.0.1:8080/article/v2/related?category_id=1'`  
`curl GET '127.0.0.1:8080/article/v2/recommend'`  
