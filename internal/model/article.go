package model

// ArticleEsAlias es 索引别名
const ArticleEsAlias = "article"

const (
	ArticleIsSolveYes = iota + 1 // 已解决
	ArticleIsSolveNo             // 未解决
)

const (
	ArticleIsRecommendYes = iota + 1 // 推荐
	ArticleIsRecommendNo             // 未推荐
)

// Article 表结构
type Article struct {
	Id          uint32 `json:"id"`           // Id
	CategoryId  uint8  `json:"category_id"`  // 分类
	Title       string `json:"title"`        // 标题
	Content     string `json:"content"`      // 内容
	BrowsNum    uint8  `json:"brows_num"`    // 浏览量
	CollectNum  uint8  `json:"collect_num"`  // 收藏量
	UpvoteNum   uint8  `json:"upvote_num"`   // 点赞量
	IsRecommend uint8  `json:"is_recommend"` // 是否推荐:1=是;2=否
	IsSolve     uint8  `json:"is_solve"`     // 是否解决:1=是;2=否
	CreatedAt   string `json:"created_at"`   // 创建时间
	UpdatedAt   string `json:"updated_at"`   // 更新时间
}

// SearchRequest 搜索请求结构
type SearchRequest struct {
	Keyword    string `form:"keyword"`                                // 关键词
	CategoryId uint8  `form:"category_id"`                            // 分类
	Sort       uint8  `form:"sort" binding:"omitempty,oneof=1 2 3"`   // 排序 1=浏览量；2=收藏；3=点赞；
	IsSolve    uint8  `form:"is_solve" binding:"omitempty,oneof=1 2"` // 是否解决
	Page       int    `form:"page,default=1"`                         // 页数
	PageSize   int    `form:"page_size,default=10"`                   // 每页数量
}

// RelatedRequest 相关文章请求结构
type RelatedRequest struct {
	CategoryId uint8 `form:"category_id" binding:"required"` // 分类
}

// SearchResponse 搜索返回结构
type SearchResponse struct {
	Id          uint32 `json:"id"`           // Id
	CategoryId  uint8  `json:"category_id"`  // 分类
	Title       string `json:"title"`        // 标题
	BrowsNum    uint8  `json:"brows_num"`    // 浏览量
	CollectNum  uint8  `json:"collect_num"`  // 收藏量
	UpvoteNum   uint8  `json:"upvote_num"`   // 点赞量
	IsRecommend uint8  `json:"is_recommend"` // 是否推荐:1=是;2=否
	IsSolve     uint8  `json:"is_solve"`     // 是否解决:1=是;2=否
	CreatedAt   string `json:"created_at"`   // 创建时间
}

// RecommendResponse 推荐文章
type RecommendResponse struct {
	Id          uint32 `json:"id"`           // Id
	Title       string `json:"title"`        // 标题
	BrowsNum    uint8  `json:"brows_num"`    // 浏览量
	CollectNum  uint8  `json:"collect_num"`  // 收藏量
	IsRecommend uint8  `json:"is_recommend"` // 是否推荐:1=是;2=否
	UpvoteNum   uint8  `json:"upvote_num"`   // 点赞量
	CreatedAt   string `json:"created_at"`   // 创建时间
}

// RelatedResponse 相关文章
type RelatedResponse struct {
	Id         uint32 `json:"id"`          // Id
	Title      string `json:"title"`       // 标题
	BrowsNum   uint8  `json:"brows_num"`   // 浏览量
	CategoryId uint8  `json:"category_id"` // 分类
}
