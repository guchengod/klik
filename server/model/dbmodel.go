package model

import "time"

// DBUser 数据库用户模型
type DBUser struct {
	ID               int       `db:"id" json:"-"`
	UID              string    `db:"uid" json:"uid"`
	Nickname         string    `db:"nickname" json:"nickname"`
	Avatar168x168URI string    `db:"avatar_168x168_uri" json:"-"`
	Avatar168x168URL string    `db:"avatar_168x168_url" json:"-"`
	Avatar300x300URI string    `db:"avatar_300x300_uri" json:"-"`
	Avatar300x300URL string    `db:"avatar_300x300_url" json:"-"`
	Gender           int       `db:"gender" json:"gender"`
	Signature        string    `db:"signature" json:"signature"`
	IPLocation       string    `db:"ip_location" json:"ip_location"`
	Province         string    `db:"province" json:"province"`
	City             string    `db:"city" json:"city"`
	Country          string    `db:"country" json:"country"`
	FollowerCount    int       `db:"follower_count" json:"follower_count"`
	FollowingCount   int       `db:"following_count" json:"following_count"`
	TotalFavorited   int       `db:"total_favorited" json:"total_favorited"`
	AwemeCount       int       `db:"aweme_count" json:"aweme_count"`
	UniqueID         string    `db:"unique_id" json:"unique_id"`
	ShortID          string    `db:"short_id" json:"short_id"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`

	// 关联字段，不在数据库表中
	Avatar168x168 AvatarInfo  `db:"-" json:"avatar_168x168"`
	Avatar300x300 AvatarInfo  `db:"-" json:"avatar_300x300"`
	CoverURL      []CoverInfo `db:"-" json:"cover_url"`
	WhiteCoverURL []CoverInfo `db:"-" json:"white_cover_url"`
}

// DBCoverURL 数据库封面URL模型
type DBCoverURL struct {
	ID        int       `db:"id" json:"-"`
	UserID    string    `db:"commenter_id" json:"-"`
	URI       string    `db:"uri" json:"-"`
	URL       string    `db:"url" json:"-"`
	Type      string    `db:"cover_type" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

// DBMusic 数据库音乐模型
type DBMusic struct {
	ID            int64     `db:"id" json:"id"`
	Title         string    `db:"title" json:"title"`
	Author        string    `db:"author" json:"author"`
	CoverURI      string    `db:"cover_uri" json:"-"`
	CoverURL      string    `db:"cover_url" json:"-"`
	PlayURL       string    `db:"play_url" json:"-"`
	Duration      int       `db:"duration" json:"duration"`
	OwnerID       string    `db:"owner_id" json:"owner_id"`
	OwnerNickname string    `db:"owner_nickname" json:"owner_nickname"`
	IsOriginal    bool      `db:"is_original" json:"is_original"`
	CreatedAt     time.Time `db:"created_at" json:"-"`

	// 关联字段，不在数据库表中
	CoverMedium CoverInfo `db:"-" json:"cover_medium"`
	CoverThumb  CoverInfo `db:"-" json:"cover_thumb"`
	PlayURLObj  CoverInfo `db:"-" json:"play_url"`
}

// DBVideo 数据库视频模型
type DBVideo struct {
	ID              int       `db:"id" json:"-"`
	AwemeID         string    `db:"aweme_id" json:"aweme_id"`
	Desc            string    `db:"video_desc" json:"desc"`
	CreateTime      int64     `db:"create_time" json:"create_time"`
	MusicID         int64     `db:"music_id" json:"-"`
	AuthorUserID    string    `db:"author_user_id" json:"author_user_id"`
	Duration        int       `db:"duration" json:"duration"`
	Type            string    `db:"video_type" json:"type"`
	ShareURL        string    `db:"share_url" json:"share_url"`
	IsTop           bool      `db:"is_top" json:"is_top"`
	PreventDownload bool      `db:"prevent_download" json:"prevent_download"`
	CreatedAt       time.Time `db:"created_at" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"-"`

	// 关联字段，不在数据库表中
	Music     DBMusic `db:"-" json:"music"`
	Author    DBUser  `db:"-" json:"author"`
	VideoInfo struct { // 视频信息
		PlayAddr       CoverInfo `json:"play_addr"`
		Cover          CoverInfo `json:"cover"`
		Height         int       `json:"height"`
		Width          int       `json:"width"`
		Ratio          string    `json:"ratio"`
		UseStaticCover bool      `json:"use_static_cover"`
		Duration       int       `json:"duration"`
	} `db:"-" json:"video"`
	StatisticsInfo struct { // 统计信息
		CommentCount int `json:"comment_count"`
		DiggCount    int `json:"digg_count"`
		CollectCount int `json:"collect_count"`
		PlayCount    int `json:"play_count"`
		ShareCount   int `json:"share_count"`
	} `db:"-" json:"statistics"`
	StatusInfo struct { // 状态信息
		IsDelete      bool `json:"is_delete"`
		AllowShare    bool `json:"allow_share"`
		IsProhibited  bool `json:"is_prohibited"`
		InReviewing   bool `json:"in_reviewing"`
		PrivateStatus int  `json:"private_status"`
	} `db:"-" json:"status"`
}

// DBVideoPlayAddress 数据库视频播放地址模型
type DBVideoPlayAddress struct {
	ID        int       `db:"id" json:"-"`
	VideoID   int       `db:"video_id" json:"-"`
	URI       string    `db:"uri" json:"uri"`
	URL       string    `db:"url" json:"-"`
	Width     int       `db:"width" json:"width"`
	Height    int       `db:"height" json:"height"`
	DataSize  int64     `db:"data_size" json:"data_size"`
	FileHash  string    `db:"file_hash" json:"file_hash"`
	CreatedAt time.Time `db:"created_at" json:"-"`

	// 关联字段，不在数据库表中
	URLList []string `db:"-" json:"url_list"`
}

// DBVideoCover 数据库视频封面模型
type DBVideoCover struct {
	ID        int       `db:"id" json:"-"`
	VideoID   int       `db:"video_id" json:"-"`
	URI       string    `db:"uri" json:"uri"`
	URL       string    `db:"url" json:"-"`
	Width     int       `db:"width" json:"width"`
	Height    int       `db:"height" json:"height"`
	CreatedAt time.Time `db:"created_at" json:"-"`

	// 关联字段，不在数据库表中
	URLList []string `db:"-" json:"url_list"`
}

// DBVideoStatistics 数据库视频统计模型
type DBVideoStatistics struct {
	ID           int       `db:"id" json:"-"`
	VideoID      int       `db:"video_id" json:"-"`
	CommentCount int       `db:"comment_count" json:"comment_count"`
	DiggCount    int       `db:"digg_count" json:"digg_count"`
	CollectCount int       `db:"collect_count" json:"collect_count"`
	PlayCount    int       `db:"play_count" json:"play_count"`
	ShareCount   int       `db:"share_count" json:"share_count"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

// DBVideoStatus 数据库视频状态模型
type DBVideoStatus struct {
	ID            int       `db:"id" json:"-"`
	VideoID       int       `db:"video_id" json:"-"`
	IsDelete      bool      `db:"is_delete" json:"is_delete"`
	AllowShare    bool      `db:"allow_share" json:"allow_share"`
	IsProhibited  bool      `db:"is_prohibited" json:"is_prohibited"`
	InReviewing   bool      `db:"in_reviewing" json:"in_reviewing"`
	PrivateStatus int       `db:"private_status" json:"private_status"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"-"`
}

// DBPost 数据库帖子模型
type DBPost struct {
	ID           int       `db:"id" json:"-"`
	PostID       string    `db:"post_id" json:"id"`
	Text         string    `db:"post_text" json:"text"`
	AuthorID     int       `db:"author_id" json:"-"`
	CreateTime   int64     `db:"create_time" json:"create_time"`
	DiggCount    int       `db:"digg_count" json:"digg_count"`
	CommentCount int       `db:"comment_count" json:"comment_count"`
	ShareCount   int       `db:"share_count" json:"share_count"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`

	// 关联字段，不在数据库表中
	Author User     `db:"-" json:"author"`
	Images []string `db:"-" json:"images,omitempty"`
}

// DBPostImage 数据库帖子图片模型
type DBPostImage struct {
	ID        int       `db:"id" json:"-"`
	PostID    int       `db:"post_id" json:"-"`
	ImageURL  string    `db:"image_url" json:"url"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

// DBGood 数据库商品模型
type DBGood struct {
	ID          int       `db:"id" json:"-"`
	GoodID      string    `db:"good_id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	Image       string    `db:"image" json:"image"`
	SaleCount   int       `db:"sale_count" json:"sale_count"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}
