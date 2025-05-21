package model


// Response 通用响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// PageResponse 分页响应
type PageResponse struct {
	PageNo int         `json:"pageNo,omitempty"`
	Total  int         `json:"total"`
	List   interface{} `json:"list"`
}

// AvatarInfo 头像信息
type AvatarInfo struct {
	URLList []string `json:"url_list"`
}

// CoverInfo 封面信息
type CoverInfo struct {
	URLList []string `json:"url_list"`
}

// Video 视频
type Video struct {
	AwemeID         string         `json:"aweme_id"`
	Desc            string         `json:"desc"`
	CreateTime      int64          `json:"create_time"`
	Music           MusicInfo      `json:"music"`
	VideoInfo       VideoInfo      `json:"video"`
	ShareURL        string         `json:"share_url"`
	Statistics      Statistics     `json:"statistics"`
	Status          StatusInfo     `json:"status"`
	TextExtra       []TextExtra    `json:"text_extra"`
	IsTop           int            `json:"is_top"`
	ShareInfo       ShareInfo      `json:"share_info"`
	Duration        int            `json:"duration"`
	ImageInfos      interface{}    `json:"image_infos"`
	RiskInfos       RiskInfo       `json:"risk_infos"`
	Position        interface{}    `json:"position"`
	AuthorUserID    string         `json:"author_user_id"`
	PreventDownload bool           `json:"prevent_download"`
	LongVideo       interface{}    `json:"long_video"`
	AwemeControl    AwemeControl   `json:"aweme_control"`
	Images          interface{}    `json:"images"`
	SuggestWords    SuggestWords   `json:"suggest_words"`
	Author          Author         `json:"author"`
}

// Comment 评论
type Comment struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	User         User      `json:"user"`
	CreateTime   int64     `json:"create_time"`
	DiggCount    int       `json:"digg_count"`
	ReplyComment []Comment `json:"reply_comment,omitempty"`
}

// Post 帖子
type Post struct {
	ID           string   `json:"id"`
	Text         string   `json:"text"`
	Author       User     `json:"author"`
	CreateTime   int64    `json:"create_time"`
	DiggCount    int      `json:"digg_count"`
	CommentCount int      `json:"comment_count"`
	ShareCount   int      `json:"share_count"`
	Images       []string `json:"images,omitempty"`
}

// Good 商品
type Good struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	SaleCount   int     `json:"sale_count"`
}

// CollectResponse 收藏响应
type CollectResponse struct {
	Video PageResponse `json:"video"`
	Music PageResponse `json:"music"`
}

// Music 音乐
type Music struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Cover    string `json:"cover"`
	PlayURL  string `json:"play_url"`
	Duration int    `json:"duration"`
}

// MusicInfo 音乐信息
type MusicInfo struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	CoverMedium   CoverMedia `json:"cover_medium"`
	CoverThumb    CoverMedia `json:"cover_thumb"`
	PlayURL       PlayMedia  `json:"play_url"`
	Duration      int        `json:"duration"`
	UserCount     int        `json:"user_count"`
	OwnerID       string     `json:"owner_id"`
	OwnerNickname string     `json:"owner_nickname"`
	IsOriginal    bool       `json:"is_original"`
}

// CoverMedia 封面媒体
type CoverMedia struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
}

// PlayMedia 播放媒体
type PlayMedia struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
	URLKey  string   `json:"url_key"`
}

// VideoInfo 视频信息
type VideoInfo struct {
	PlayAddr      PlayAddr `json:"play_addr"`
	Cover         Cover    `json:"cover"`
	Height        int      `json:"height"`
	Width         int      `json:"width"`
	Ratio         string   `json:"ratio"`
	UseStaticCover bool     `json:"use_static_cover"`
	Duration      int      `json:"duration"`
}

// PlayAddr 播放地址
type PlayAddr struct {
	URI      string   `json:"uri"`
	URLList  []string `json:"url_list"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	URLKey   string   `json:"url_key"`
	DataSize int64    `json:"data_size"`
	FileHash string   `json:"file_hash"`
	FileCS   string   `json:"file_cs"`
}

// Cover 封面
type Cover struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
}

// Statistics 统计信息
type Statistics struct {
	AdmireCount  int `json:"admire_count"`
	CommentCount int `json:"comment_count"`
	DiggCount    int `json:"digg_count"`
	CollectCount int `json:"collect_count"`
	PlayCount    int `json:"play_count"`
	ShareCount   int `json:"share_count"`
}

// StatusInfo 状态信息
type StatusInfo struct {
	ListenVideoStatus int          `json:"listen_video_status"`
	IsDelete         bool         `json:"is_delete"`
	AllowShare       bool         `json:"allow_share"`
	IsProhibited     bool         `json:"is_prohibited"`
	InReviewing      bool         `json:"in_reviewing"`
	PartSee          int          `json:"part_see"`
	PrivateStatus    int          `json:"private_status"`
	ReviewResult     ReviewResult `json:"review_result"`
}

// ReviewResult 审核结果
type ReviewResult struct {
	ReviewStatus int `json:"review_status"`
}

// TextExtra 文本额外信息
type TextExtra struct {
	Start        int    `json:"start"`
	End          int    `json:"end"`
	Type         int    `json:"type"`
	HashtagName  string `json:"hashtag_name"`
	HashtagID    string `json:"hashtag_id"`
	IsCommerce   bool   `json:"is_commerce"`
	CaptionStart int    `json:"caption_start"`
	CaptionEnd   int    `json:"caption_end"`
}

// ShareInfo 分享信息
type ShareInfo struct {
	ShareURL      string `json:"share_url"`
	ShareLinkDesc string `json:"share_link_desc"`
}

// RiskInfo 风险信息
type RiskInfo struct {
	Vote     bool   `json:"vote"`
	Warn     bool   `json:"warn"`
	RiskSink bool   `json:"risk_sink"`
	Type     int    `json:"type"`
	Content  string `json:"content"`
}

// AwemeControl 控制信息
type AwemeControl struct {
	CanForward    bool `json:"can_forward"`
	CanShare      bool `json:"can_share"`
	CanComment    bool `json:"can_comment"`
	CanShowComment bool `json:"can_show_comment"`
}

// SuggestWords 建议词
type SuggestWords struct {
	SuggestWords []SuggestWord `json:"suggest_words"`
}

// SuggestWord 建议词
type SuggestWord struct {
	Words     []Word  `json:"words"`
	Scene     string  `json:"scene"`
	IconURL   string  `json:"icon_url"`
	HintText  string  `json:"hint_text"`
	ExtraInfo string  `json:"extra_info"`
}

// Word 词
type Word struct {
	Word   string `json:"word"`
	WordID string `json:"word_id"`
	Info   string `json:"info"`
}

// Author 作者
type Author struct {
	Avatar168x168   Avatar      `json:"avatar_168x168"`
	Avatar300x300   Avatar      `json:"avatar_300x300"`
	AwemeCount      int         `json:"aweme_count"`
	BirthdayHideLevel int        `json:"birthday_hide_level"`
	CanShowGroupCard int         `json:"can_show_group_card"`
	CardEntries      []CardEntry `json:"card_entries"`
	City             string      `json:"city"`
	CommerceInfo    CommerceInfo `json:"commerce_info"`
	CommerceUserInfo CommerceUserInfo `json:"commerce_user_info"`
	CommerceUserLevel int        `json:"commerce_user_level"`
	Country          string      `json:"country"`
	CoverColour      string      `json:"cover_colour"`
	CoverURL         []CoverURL  `json:"cover_url"`
	District         string      `json:"district"`
	FavoritingCount  int         `json:"favoriting_count"`
	FollowStatus     int         `json:"follow_status"`
	FollowerCount    int         `json:"follower_count"`
	FollowerRequestStatus int    `json:"follower_request_status"`
	FollowerStatus    int        `json:"follower_status"`
	FollowingCount    int        `json:"following_count"`
	ForwardCount      int        `json:"forward_count"`
	Gender            int        `json:"gender"`
	IPLocation        string     `json:"ip_location"`
	MaxFollowerCount  int        `json:"max_follower_count"`
	MplatformFollowersCount int `json:"mplatform_followers_count"`
	Nickname          string     `json:"nickname"`
	Province          string     `json:"province"`
	PublicCollectsCount int      `json:"public_collects_count"`
	ShareInfo         ShareInfoUser `json:"share_info"`
	ShortID           string     `json:"short_id"`
	Signature         string     `json:"signature"`
	TotalFavorited    int        `json:"total_favorited"`
	UID               string     `json:"uid"`
	UniqueID          string     `json:"unique_id"`
	UserAge           int        `json:"user_age"`
	WhiteCoverURL     []CoverURL `json:"white_cover_url"`
}

// Avatar 头像
type Avatar struct {
	Height  int      `json:"height"`
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
}

// CardEntry 卡片条目
type CardEntry struct {
	CardData     string    `json:"card_data"`
	EventParams  string    `json:"event_params"`
	GotoURL      string    `json:"goto_url"`
	IconDark     IconMedia `json:"icon_dark"`
	IconLight    IconMedia `json:"icon_light"`
	SubTitle     string    `json:"sub_title"`
	Title        string    `json:"title"`
	Type         int       `json:"type"`
}

// IconMedia 图标媒体
type IconMedia struct {
	URLList []string `json:"url_list"`
}

// CommerceInfo 商业信息
type CommerceInfo struct {
	ChallengeList    interface{}   `json:"challenge_list"`
	HeadImageList    interface{}   `json:"head_image_list"`
	OfflineInfoList  []interface{} `json:"offline_info_list"`
	SmartPhoneList   interface{}   `json:"smart_phone_list"`
	TaskList         interface{}   `json:"task_list"`
}

// CommerceUserInfo 商业用户信息
type CommerceUserInfo struct {
	AdRevenueRits           interface{} `json:"ad_revenue_rits"`
	HasAdsEntry             bool        `json:"has_ads_entry"`
	ShowStarAtlasCooperation bool        `json:"show_star_atlas_cooperation"`
	StarAtlas                int         `json:"star_atlas"`
}

// CoverURL 封面URL
type CoverURL struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

// ShareInfoUser 用户分享信息
type ShareInfoUser struct {
	BoolPersist      int          `json:"bool_persist"`
	ShareDesc        string       `json:"share_desc"`
	ShareImageURL    ShareImageURL `json:"share_image_url"`
	ShareQrcodeURL   ShareQrcodeURL `json:"share_qrcode_url"`
	ShareTitle       string       `json:"share_title"`
	ShareURL         string       `json:"share_url"`
	ShareWeiboDesc   string       `json:"share_weibo_desc"`
}

// ShareImageURL 分享图片URL
type ShareImageURL struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

// ShareQrcodeURL 分享二维码URL
type ShareQrcodeURL struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

// PageParams 分页参数
type PageParams struct {
	PageNo   int `form:"pageNo" json:"pageNo"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

// VideoRecommendParams 视频推荐参数
type VideoRecommendParams struct {
	Start    int `form:"start" json:"start"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

// GetPageRange 获取页面范围
func GetPageRange(params PageParams) (int, int) {
	offset := params.PageNo * params.PageSize
	limit := params.PageNo*params.PageSize + params.PageSize
	return offset, limit
}

// User 用户
type User struct {
	UID           string      `json:"uid"`
	Nickname      string      `json:"nickname"`
	Gender        int         `json:"gender"`
	Signature     string      `json:"signature"`
	Avatar168x168 AvatarInfo  `json:"avatar_168x168"`
	Avatar300x300 AvatarInfo  `json:"avatar_300x300"`
	CoverURL      []CoverInfo `json:"cover_url"`
	WhiteCoverURL []CoverInfo `json:"white_cover_url"`
	UniqueID      interface{} `json:"unique_id"`
}

