package model

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
