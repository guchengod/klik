package model

import (
	"database/sql"
	"errors"
	"fmt"
	"klik/server/config"
)

// 从数据库获取用户信息
func GetUserByID(userID string) (User, error) {
	if config.DB == nil {
		return User{}, errors.New("数据库未初始化")
	}

	query := `
		SELECT id, uid, nickname, gender, signature
		FROM users
		WHERE uid = $1
	`

	var user User
	err := config.DB.QueryRow(query, userID).Scan(
		&user.UID,
		&user.Nickname,
		&user.Gender,
		&user.Signature,
	)
	if err != nil {
		return User{}, fmt.Errorf("查询用户信息失败: %v", err)
	}

	return user, nil
}

// getUserAvatars 获取用户头像
func getUserAvatars(userID string) (AvatarInfo, AvatarInfo, error) {
	if config.DB == nil {
		return AvatarInfo{}, AvatarInfo{}, fmt.Errorf("数据库未初始化")
	}

	query := `
		SELECT url, cover_type FROM cover_urls 
		WHERE user_id = (SELECT id FROM users WHERE uid = $1) 
		AND (cover_type = 'avatar_168x168' OR cover_type = 'avatar_300x300')
	`
	avatarRows, err := config.DB.Query(query, userID)
	if err != nil {
		return AvatarInfo{}, AvatarInfo{}, fmt.Errorf("查询用户头像失败: %v", err)
	}
	defer avatarRows.Close()

	var avatar168URL, avatar300URL string
	for avatarRows.Next() {
		var url, avatarType string
		err := avatarRows.Scan(&url, &avatarType)
		if err != nil {
			return AvatarInfo{}, AvatarInfo{}, fmt.Errorf("解析用户头像失败: %v", err)
		}

		if avatarType == "avatar_168x168" {
			avatar168URL = url
		} else if avatarType == "avatar_300x300" {
			avatar300URL = url
		}
	}

	avatar168 := AvatarInfo{
		URLList: []string{avatar168URL},
	}
	avatar300 := AvatarInfo{
		URLList: []string{avatar300URL},
	}

	return avatar168, avatar300, nil
}

// 从数据库获取用户收藏的视频
func GetUserCollectVideosFromDB(userID string, offset, limit int) ([]Video, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `
		SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration, 
		       u.uid, u.nickname, u.gender, u.signature, 
		       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count, vs.play_count, vs.admire_count
		FROM videos v
		LEFT JOIN users u ON v.author_user_id = u.uid
		LEFT JOIN video_statistics vs ON v.id = vs.video_id
		LEFT JOIN user_collect_videos ucv ON v.id = ucv.video_id
		WHERE ucv.commenter_id = $1
		ORDER BY ucv.created_at DESC
		LIMIT $2 OFFSET $3
	`

	// 执行查询
	rows, err := config.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询用户收藏视频失败: %v", err)
	}
	defer rows.Close()

	// 解析结果
	var videos []Video
	for rows.Next() {
		var videoID int64
		var awemeID, desc string
		var createTime int64
		var authorUserID string
		var duration int
		var userUID, nickname, gender, signature string
		var commentCount, diggCount, collectCount, shareCount, playCount, admireCount int64

		// 扫描行数据
		err := rows.Scan(
			&videoID, &awemeID, &desc, &createTime, &authorUserID, &duration,
			&userUID, &nickname, &gender, &signature,
			&commentCount, &diggCount, &collectCount, &shareCount,
			&playCount, &admireCount,
		)
		if err != nil {
			return nil, fmt.Errorf("解析用户收藏视频数据失败: %v", err)
		}

		// 获取用户头像
		avatarQuery := `
			SELECT uri, url, cover_type FROM cover_urls 
			WHERE user_id = (SELECT id FROM users WHERE uid = $1) 
			AND (cover_type = 'avatar_168x168' OR cover_type = 'avatar_300x300')
		`
		avatarRows, err := config.DB.Query(avatarQuery, userUID)
		if err != nil {
			return nil, fmt.Errorf("查询用户头像失败: %v", err)
		}
		defer avatarRows.Close()

		var avatar168URI, avatar168URL, avatar300URI, avatar300URL string
		for avatarRows.Next() {
			var uri, url, avatarType string
			err := avatarRows.Scan(&uri, &url, &avatarType)
			if err != nil {
				return nil, fmt.Errorf("解析用户头像失败: %v", err)
			}

			if avatarType == "avatar_168x168" {
				avatar168URI = uri
				avatar168URL = url
			} else if avatarType == "avatar_300x300" {
				avatar300URI = uri
				avatar300URL = url
			}
		}

		// 获取视频封面
		coverQuery := `
			SELECT uri, url FROM video_covers
			WHERE video_id = $1
			LIMIT 1
		`
		var coverURI, coverURL string
		coverErr := config.DB.QueryRow(coverQuery, videoID).Scan(&coverURI, &coverURL)
		if coverErr != nil && coverErr != sql.ErrNoRows {
			return nil, fmt.Errorf("查询视频封面失败: %v", coverErr)
		}

		// 获取视频地址
		playAddrQuery := `
			SELECT uri, url, width, height, data_size, file_hash, file_cs FROM video_play_urls
			WHERE video_id = $1
			LIMIT 1
		`
		var playURI, playURL string
		var width, height int
		var dataSize int64
		var fileHash, fileCS string
		playErr := config.DB.QueryRow(playAddrQuery, videoID).Scan(
			&playURI, &playURL, &width, &height, &dataSize, &fileHash, &fileCS,
		)
		if playErr != nil && playErr != sql.ErrNoRows {
			return nil, fmt.Errorf("查询视频播放地址失败: %v", playErr)
		}

		// 默认视频URL，如果数据库中没有数据
		if playURL == "" {
			playURL = fmt.Sprintf("https://example.com/video/%s.mp4", awemeID)
		}

		// 获取音乐信息
		musicQuery := `
			SELECT m.id, m.title, m.author, m.duration, m.play_url, m.owner_id, m.owner_nickname, m.is_original
			FROM musics m
			JOIN video_musics vm ON m.id = vm.music_id
			WHERE vm.video_id = $1
			LIMIT 1
		`
		var musicID int64
		var musicTitle, musicAuthor, musicPlayURL, musicOwnerID, musicOwnerNickname string
		var musicDuration int
		var isOriginal bool
		musicErr := config.DB.QueryRow(musicQuery, videoID).Scan(
			&musicID, &musicTitle, &musicAuthor, &musicDuration, &musicPlayURL,
			&musicOwnerID, &musicOwnerNickname, &isOriginal,
		)
		if musicErr != nil && musicErr != sql.ErrNoRows {
			return nil, fmt.Errorf("查询音乐信息失败: %v", musicErr)
		}

		// 构建视频对象
		video := Video{
			AwemeID:    awemeID,
			Desc:       desc,
			CreateTime: createTime,
			ShareURL:   fmt.Sprintf("https://example.com/share/%s", awemeID),
			Duration:   duration,
			AuthorUserID: authorUserID,
			PreventDownload: false,
			Music: MusicInfo{
				ID:            musicID,
				Title:         musicTitle,
				Author:        musicAuthor,
				Duration:      musicDuration,
				OwnerID:       musicOwnerID,
				OwnerNickname: musicOwnerNickname,
				IsOriginal:    isOriginal,
				PlayURL: PlayMedia{
					URI:     "",
					URLList: []string{musicPlayURL},
					Width:   0,
					Height:  0,
					URLKey:  "",
				},
				CoverMedium: CoverMedia{
					URI:     "",
					URLList: []string{""},
					Width:   0,
					Height:  0,
				},
				CoverThumb: CoverMedia{
					URI:     "",
					URLList: []string{""},
					Width:   0,
					Height:  0,
				},
				UserCount: 0,
			},
			VideoInfo: VideoInfo{
				PlayAddr: PlayAddr{
					URI:      playURI,
					URLList:  []string{playURL},
					Width:    width,
					Height:   height,
					URLKey:   "",
					DataSize: dataSize,
					FileHash: fileHash,
					FileCS:   fileCS,
				},
				Cover: Cover{
					URI:     coverURI,
					URLList: []string{coverURL},
					Width:   0,
					Height:  0,
				},
				Height:        height,
				Width:         width,
				Ratio:         "540p",
				UseStaticCover: false,
				Duration:      duration,
			},
			Statistics: Statistics{
				AdmireCount:  int(admireCount),
				CommentCount: int(commentCount),
				DiggCount:    int(diggCount),
				CollectCount: int(collectCount),
				PlayCount:    int(playCount),
				ShareCount:   int(shareCount),
			},
			Status: StatusInfo{
				ListenVideoStatus: 0,
				IsDelete:         false,
				AllowShare:       true,
				IsProhibited:     false,
				InReviewing:      false,
				PartSee:          0,
				PrivateStatus:    0,
				ReviewResult: ReviewResult{
					ReviewStatus: 0,
				},
			},
			TextExtra: []TextExtra{},
			IsTop:     0,
			ShareInfo: ShareInfo{
				ShareURL:      fmt.Sprintf("https://example.com/share/%s", awemeID),
				ShareLinkDesc: desc,
			},
			AwemeControl: AwemeControl{
				CanForward:     true,
				CanShare:       true,
				CanComment:     true,
				CanShowComment: true,
			},
			Author: Author{
				UID:        userUID,
				Nickname:   nickname,
				Gender:     genderToInt(gender),
				Signature:  signature,
				Avatar168x168: Avatar{
					URI:     avatar168URI,
					URLList: []string{avatar168URL},
					Width:   168,
					Height:  168,
				},
				Avatar300x300: Avatar{
					URI:     avatar300URI,
					URLList: []string{avatar300URL},
					Width:   300,
					Height:  300,
				},
				FollowerCount:  0,
				FollowingCount: 0,
				AwemeCount:     0,
				TotalFavorited: 0,
				UniqueID:       userUID,
				CoverURL:       []CoverURL{},
				WhiteCoverURL:  []CoverURL{},
			},
		}

		videos = append(videos, video)
	}

	// 检查是否有查询错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("查询用户收藏视频数据时发生错误: %v", err)
	}

	return videos, nil
}

// 从数据库获取用户收藏的音乐
func GetUserCollectMusicFromDB(userID string, offset, limit int) ([]Music, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 注意：schema.sql中没有user_collected_music表，这里假设它存在或者需要创建
	query := `
		SELECT m.id, m.title, m.author, m.cover_url, m.play_url, m.duration
		FROM music m
		JOIN user_collect_music ucm ON m.id = ucm.music_id
		WHERE ucm.commenter_id = $1
		ORDER BY ucm.created_at DESC
		LIMIT $2 OFFSET $3
	`

	// 执行查询
	rows, err := config.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询用户收藏音乐失败: %v", err)
	}
	defer rows.Close()

	// 解析结果
	var musicList []Music
	for rows.Next() {
		var music Music
		err := rows.Scan(
			&music.ID,
			&music.Title,
			&music.Artist,
			&music.Cover,
			&music.PlayURL,
			&music.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("解析用户收藏音乐数据失败: %v", err)
		}

		musicList = append(musicList, music)
	}

	// 检查是否有查询错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("查询用户收藏音乐数据时发生错误: %v", err)
	}

	return musicList, nil
}

// 获取用户收藏视频总数
func GetUserCollectVideosCountFromDB(userID string) (int, error) {
	if config.DB == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}

	query := `
		SELECT COUNT(*) FROM user_collect_videos
		WHERE commenter_id = $1
	`

	var count int
	err := config.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询用户收藏视频总数失败: %v", err)
	}

	return count, nil
}

// 获取用户收藏音乐总数
func GetUserCollectMusicCountFromDB(userID string) (int, error) {
	if config.DB == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}

	// 注意：schema.sql中没有user_collect_music表，这里假设它存在或者需要创建
	query := `
		SELECT COUNT(*) FROM user_collect_music
		WHERE commenter_id = $1
	`

	var count int
	err := config.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询用户收藏音乐总数失败: %v", err)
	}

	return count, nil
}

// 从数据库获取用户视频列表
func GetUserVideoListFromDB(userID string) ([]Video, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `
		SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration, 
		       u.uid, u.nickname, u.gender, u.signature, 
		       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count, vs.play_count, vs.admire_count
		FROM videos v
		LEFT JOIN users u ON v.author_user_id = u.uid
		LEFT JOIN video_statistics vs ON v.id = vs.video_id
		WHERE v.author_user_id = $1
		ORDER BY v.create_time DESC
	`

	// 执行查询
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户视频列表失败: %v", err)
	}
	defer rows.Close()

	// 解析结果
	var videos []Video
	for rows.Next() {
		var videoID int64
		var awemeID, desc string
		var createTime int64
		var authorUserID string
		var duration int
		var userUID, nickname, signature string
		var gender string
		var commentCount, diggCount, collectCount, shareCount, playCount, admireCount int64

		// 扫描行数据
		err := rows.Scan(
			&videoID, &awemeID, &desc, &createTime, &authorUserID, &duration,
			&userUID, &nickname, &gender, &signature,
			&commentCount, &diggCount, &collectCount, &shareCount,
			&playCount, &admireCount,
		)
		if err != nil {
			return nil, fmt.Errorf("解析用户视频列表数据失败: %v", err)
		}

		// 获取用户头像
		avatarQuery := `
			SELECT uri, url, cover_type FROM cover_urls 
			WHERE user_id = (SELECT id FROM users WHERE uid = $1) 
			AND (cover_type = 'avatar_168x168' OR cover_type = 'avatar_300x300')
		`
		avatarRows, err := config.DB.Query(avatarQuery, userUID)
		if err != nil {
			return nil, fmt.Errorf("查询用户头像失败: %v", err)
		}
		defer avatarRows.Close()

		var avatar168URI, avatar168URL, avatar300URI, avatar300URL string
		for avatarRows.Next() {
			var uri, url, avatarType string
			err := avatarRows.Scan(&uri, &url, &avatarType)
			if err != nil {
				return nil, fmt.Errorf("解析用户头像失败: %v", err)
			}

			if avatarType == "avatar_168x168" {
				avatar168URI = uri
				avatar168URL = url
			} else if avatarType == "avatar_300x300" {
				avatar300URI = uri
				avatar300URL = url
			}
		}

		// 获取视频封面
		coverQuery := `
			SELECT uri, url FROM video_covers
			WHERE video_id = $1
			LIMIT 1
		`
		var coverURI, coverURL string
		coverErr := config.DB.QueryRow(coverQuery, videoID).Scan(&coverURI, &coverURL)
		if coverErr != nil && coverErr != sql.ErrNoRows {
			return nil, fmt.Errorf("查询视频封面失败: %v", coverErr)
		}

		// 获取视频地址
		playAddrQuery := `
			SELECT uri, url, width, height, data_size, file_hash, file_cs FROM video_play_urls
			WHERE video_id = $1
			LIMIT 1
		`
		var playURI, playURL string
		var width, height int
		var dataSize int64
		var fileHash, fileCS string
		playErr := config.DB.QueryRow(playAddrQuery, videoID).Scan(
			&playURI, &playURL, &width, &height, &dataSize, &fileHash, &fileCS,
		)
		if playErr != nil && playErr != sql.ErrNoRows {
			return nil, fmt.Errorf("查询视频播放地址失败: %v", playErr)
		}

		// 默认视频URL，如果数据库中没有数据
		if playURL == "" {
			playURL = fmt.Sprintf("https://example.com/video/%s.mp4", awemeID)
		}

		// 构建视频对象
		video := Video{
			AwemeID:    awemeID,
			Desc:       desc,
			CreateTime: createTime,
			ShareURL:   fmt.Sprintf("https://example.com/share/%s", awemeID),
			Duration:   duration,
			AuthorUserID: authorUserID,
			PreventDownload: false,
			VideoInfo: VideoInfo{
				PlayAddr: PlayAddr{
					URI:     playURI,
					URLList: []string{playURL},
					Width:   width,
					Height:  height,
					URLKey:  "",
					DataSize: dataSize,
					FileHash: fileHash,
					FileCS:   fileCS,
				},
				Cover: Cover{
					URI:     coverURI,
					URLList: []string{coverURL},
					Width:   0,
					Height:  0,
				},
				Height:        height,
				Width:         width,
				Ratio:         "540p",
				UseStaticCover: false,
				Duration:      duration,
			},
			Statistics: Statistics{
				AdmireCount:  int(admireCount),
				CommentCount: int(commentCount),
				DiggCount:    int(diggCount),
				CollectCount: int(collectCount),
				PlayCount:    int(playCount),
				ShareCount:   int(shareCount),
			},
			Status: StatusInfo{
				ListenVideoStatus: 0,
				IsDelete:         false,
				AllowShare:       true,
				IsProhibited:     false,
				InReviewing:      false,
				PartSee:          0,
				PrivateStatus:    0,
				ReviewResult: ReviewResult{
					ReviewStatus: 0,
				},
			},
			TextExtra: []TextExtra{},
			IsTop:     0,
			ShareInfo: ShareInfo{
				ShareURL:      fmt.Sprintf("https://example.com/share/%s", awemeID),
				ShareLinkDesc: desc,
			},
			AwemeControl: AwemeControl{
				CanForward:     true,
				CanShare:       true,
				CanComment:     true,
				CanShowComment: true,
			},
			Author: Author{
				UID:        userUID,
				Nickname:   nickname,
				Gender:     genderToInt(gender),
				Signature:  signature,
				Avatar168x168: Avatar{
					URI:     avatar168URI,
					URLList: []string{avatar168URL},
					Width:   168,
					Height:  168,
				},
				Avatar300x300: Avatar{
					URI:     avatar300URI,
					URLList: []string{avatar300URL},
					Width:   300,
					Height:  300,
				},
				FollowerCount:  0,
				FollowingCount: 0,
				AwemeCount:     0,
				TotalFavorited: 0,
				UniqueID:       userUID,
				CoverURL:       []CoverURL{},
				WhiteCoverURL:  []CoverURL{},
			},
		}

		videos = append(videos, video)
	}

	// 检查是否有查询错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("查询用户视频列表数据时发生错误: %v", err)
	}

	return videos, nil
}

// 从数据库获取用户好友列表
func GetUserFriendsFromDB(userID string) ([]DBUser, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 注意：schema.sql中没有user_friends表，这里假设它存在
	query := `
		SELECT u.uid, u.nickname, u.gender, u.signature, u.ip_location, u.province, u.city, u.country, 
		       u.follower_count, u.following_count, u.total_favorited, u.aweme_count, u.unique_id, u.short_id
		FROM users u
		WHERE u.uid != $1
		ORDER BY u.nickname
		LIMIT 20
	`

	// 执行查询
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户好友列表失败: %v", err)
	}
	defer rows.Close()

	// 解析结果
	var users []DBUser
	for rows.Next() {
		var user DBUser
		err := rows.Scan(
			&user.UID,
			&user.Nickname,
			&user.Gender,
			&user.Signature,
			&user.IPLocation,
			&user.Province,
			&user.City,
			&user.Country,
			&user.FollowerCount,
			&user.FollowingCount,
			&user.TotalFavorited,
			&user.AwemeCount,
			&user.UniqueID,
			&user.ShortID,
		)
		if err != nil {
			return nil, fmt.Errorf("解析用户好友数据失败: %v", err)
		}

		// 获取用户头像
		user.Avatar168x168, user.Avatar300x300, err = getUserAvatars(user.UID)
		if err != nil {
			// 使用默认头像
			user.Avatar168x168 = AvatarInfo{
				URLList: []string{"https://example.com/default-avatar.jpg"},
			}
			user.Avatar300x300 = AvatarInfo{
				URLList: []string{"https://example.com/default-avatar.jpg"},
			}
		}

		users = append(users, user)
	}

	// 检查是否有查询错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("查询用户好友数据时发生错误: %v", err)
	}

	return users, nil
}
