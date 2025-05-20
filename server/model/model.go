package model

// 通用响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 分页响应
type PageResponse struct {
	PageNo int         `json:"pageNo,omitempty"`
	Total  int         `json:"total"`
	List   interface{} `json:"list"`
}

// 用户
type User struct {
	UID           string       `json:"uid"`
	Nickname      string       `json:"nickname"`
	Avatar168x168 AvatarInfo   `json:"avatar_168x168"`
	Avatar300x300 AvatarInfo   `json:"avatar_300x300"`
	CoverURL      []CoverInfo  `json:"cover_url"`
	WhiteCoverURL []CoverInfo  `json:"white_cover_url"`
	UniqueID      interface{}  `json:"unique_id"`
}

// 头像信息
type AvatarInfo struct {
	URLList []string `json:"url_list"`
}

// 封面信息
type CoverInfo struct {
	URLList []string `json:"url_list"`
}

// 视频
type Video struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Src           string      `json:"src"`
	Author        User        `json:"author"`
	AuthorUserID  string      `json:"author_user_id,omitempty"`
	Description   string      `json:"description,omitempty"`
	CommentCount  int         `json:"comment_count,omitempty"`
	DiggCount     int         `json:"digg_count,omitempty"`
	ShareCount    int         `json:"share_count,omitempty"`
	CollectCount  int         `json:"collect_count,omitempty"`
}

// 评论
type Comment struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	User         User      `json:"user"`
	CreateTime   int64     `json:"create_time"`
	DiggCount    int       `json:"digg_count"`
	ReplyComment []Comment `json:"reply_comment,omitempty"`
}

// 帖子
type Post struct {
	ID          string   `json:"id"`
	Text        string   `json:"text"`
	Author      User     `json:"author"`
	CreateTime  int64    `json:"create_time"`
	DiggCount   int      `json:"digg_count"`
	CommentCount int     `json:"comment_count"`
	ShareCount  int      `json:"share_count"`
	Images      []string `json:"images,omitempty"`
}

// 商品
type Good struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	SaleCount   int     `json:"sale_count"`
}

// 收藏响应
type CollectResponse struct {
	Video PageResponse `json:"video"`
	Music PageResponse `json:"music"`
}

// 音乐
type Music struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Cover  string `json:"cover"`
}

// 分页参数
type PageParams struct {
	PageNo   int `form:"pageNo" json:"pageNo"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

// 视频推荐参数
type VideoRecommendParams struct {
	Start    int `form:"start" json:"start"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

// 获取页面范围
func GetPageRange(params PageParams) (int, int) {
	offset := params.PageNo * params.PageSize
	limit := params.PageNo*params.PageSize + params.PageSize
	return offset, limit
}
