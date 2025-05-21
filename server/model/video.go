package model

import (
	"database/sql"
	"fmt"
	"klik/server/config"
)

// 获取推荐视频列表
func GetRecommendVideosFromDB(start, pageSize int) ([]Video, error) {
	if config.UseDB && config.DB != nil {
		// 从 PostgreSQL 数据库中获取视频数据
		query := `
			SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration,
			       u.uid, u.nickname, u.gender, u.signature, 
			       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count,
			       vs.play_count, vs.collect_count
			FROM videos v
			LEFT JOIN users u ON v.author_user_id = u.uid
			LEFT JOIN video_statistics vs ON v.id = vs.video_id
			WHERE v.video_type = 'recommend-video'
			ORDER BY v.create_time DESC
			LIMIT $1 OFFSET $2
		`

		// 执行查询
		rows, err := config.DB.Query(query, pageSize, start)
		if err != nil {
			return nil, fmt.Errorf("查询视频数据失败: %v", err)
		}
		defer rows.Close()

		// 解析结果
		var videos []Video
		for rows.Next() {
			var (
				videoID, awemeID, desc, authorUserID              string
				createTime                                        int64
				duration                                          int
				userUID, nickname, signature                      string
				gender                                            int
				commentCount, diggCount, collectCount, shareCount int
				playCount, admireCount                            int
			)

			// 扫描行数据
			err := rows.Scan(
				&videoID, &awemeID, &desc, &createTime, &authorUserID, &duration,
				&userUID, &nickname, &gender, &signature,
				&commentCount, &diggCount, &collectCount, &shareCount,
				&playCount, &admireCount,
			)
			if err != nil {
				return nil, fmt.Errorf("解析视频数据失败: %v", err)
			}

			// 获取用户头像
			avatarQuery := `
				SELECT uri_path, url_path, cover_type FROM cover_urls 
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
				SELECT uri, url, width, height, data_size, file_hash FROM video_play_addresses
				WHERE video_id = $1
				LIMIT 1
			`
			var playURI, playURL string
			var width, height int
			var dataSize int64
			var fileHash, fileCS string
			playErr := config.DB.QueryRow(playAddrQuery, videoID).Scan(
				&playURI, &playURL, &width, &height, &dataSize, &fileHash,
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
				FROM music m
				JOIN videos vm ON m.id = vm.music_id
				WHERE vm.id = $1
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
				AwemeID:         awemeID,
				Desc:            desc,
				CreateTime:      createTime,
				ShareURL:        fmt.Sprintf("https://example.com/share/%s", awemeID),
				Duration:        duration,
				AuthorUserID:    authorUserID,
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
					Height:         height,
					Width:          width,
					Ratio:          "540p",
					UseStaticCover: false,
					Duration:       duration,
				},
				Statistics: Statistics{
					AdmireCount:  admireCount,
					CommentCount: commentCount,
					DiggCount:    diggCount,
					CollectCount: collectCount,
					PlayCount:    playCount,
					ShareCount:   shareCount,
				},
				Status: StatusInfo{
					ListenVideoStatus: 0,
					IsDelete:          false,
					AllowShare:        true,
					IsProhibited:      false,
					InReviewing:       false,
					PartSee:           0,
					PrivateStatus:     0,
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
					UID:       userUID,
					Nickname:  nickname,
					Gender:    gender,
					Signature: signature,
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
			return nil, fmt.Errorf("查询视频数据时发生错误: %v", err)
		}

		return videos, nil
	}

	// 如果没有使用数据库，返回错误
	return nil, fmt.Errorf("数据库未初始化")
}

// 获取长视频推荐列表
func GetLongRecommendVideosFromDB(offset, limit int) ([]Video, error) {
	if config.DB != nil {
		// 从 PostgreSQL 数据库中获取长视频数据
		query := `
			SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration,
			       u.uid, u.nickname, u.gender, u.signature, 
			       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count,
			       vs.play_count, vs.collect_count
			FROM videos v
			LEFT JOIN users u ON v.author_user_id = u.uid
			LEFT JOIN video_statistics vs ON v.id = vs.video_id
			WHERE v.video_type = 'long-video' AND v.duration > 60
			ORDER BY v.create_time DESC
			LIMIT $1 OFFSET $2
		`

		// 执行查询
		rows, err := config.DB.Query(query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("查询长视频数据失败: %v", err)
		}
		defer rows.Close()

		// 解析结果
		var videos []Video
		for rows.Next() {
			var (
				videoID, awemeID, desc, authorUserID              string
				createTime                                        int64
				duration                                          int
				userUID, nickname, signature                      string
				gender                                            int
				commentCount, diggCount, collectCount, shareCount int
				playCount, admireCount                            int
			)

			// 扫描行数据
			err := rows.Scan(
				&videoID, &awemeID, &desc, &createTime, &authorUserID, &duration,
				&userUID, &nickname, &gender, &signature,
				&commentCount, &diggCount, &collectCount, &shareCount,
				&playCount, &admireCount,
			)
			if err != nil {
				return nil, fmt.Errorf("解析长视频数据失败: %v", err)
			}

			// 获取用户头像
			avatarQuery := `
				SELECT uri_path, url_path, cover_type FROM cover_urls 
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
				SELECT uri, url, width, height, data_size, file_hash FROM video_play_addresses
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
				FROM music m
				JOIN videos vm ON m.id = vm.music_id
				WHERE vm.id = $1
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
				AwemeID:         awemeID,
				Desc:            desc,
				CreateTime:      createTime,
				ShareURL:        fmt.Sprintf("https://example.com/share/%s", awemeID),
				Duration:        duration,
				AuthorUserID:    authorUserID,
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
					Height:         height,
					Width:          width,
					Ratio:          "540p",
					UseStaticCover: false,
					Duration:       duration,
				},
				Statistics: Statistics{
					AdmireCount:  admireCount,
					CommentCount: commentCount,
					DiggCount:    diggCount,
					CollectCount: collectCount,
					PlayCount:    playCount,
					ShareCount:   shareCount,
				},
				Status: StatusInfo{
					ListenVideoStatus: 0,
					IsDelete:          false,
					AllowShare:        true,
					IsProhibited:      false,
					InReviewing:       false,
					PartSee:           0,
					PrivateStatus:     0,
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
					UID:       userUID,
					Nickname:  nickname,
					Gender:    gender,
					Signature: signature,
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
			return nil, fmt.Errorf("查询长视频数据时发生错误: %v", err)
		}

		return videos, nil
	}

	// 如果没有使用数据库，返回错误
	return nil, fmt.Errorf("数据库未初始化")
}

// 获取私有视频列表
func GetPrivateVideosFromDB(offset, limit int) ([]Video, error) {
	if config.DB != nil {
		// 从 PostgreSQL 数据库中获取私有视频数据
		query := `
			SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration,
			       u.uid, u.nickname, u.gender, u.signature, 
			       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count,
			       vs.play_count, vs.collect_count
			FROM videos v
			LEFT JOIN users u ON v.author_user_id = u.uid
			LEFT JOIN video_statistics vs ON v.id = vs.video_id
			WHERE v.video_type = 'private-video'
			ORDER BY v.create_time DESC
			LIMIT $1 OFFSET $2
		`

		// 执行查询
		rows, err := config.DB.Query(query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("查询私有视频数据失败: %v", err)
		}
		defer rows.Close()

		// 解析结果
		var videos []Video
		for rows.Next() {
			var (
				videoID, awemeID, desc, authorUserID              string
				createTime                                        int64
				duration                                          int
				userUID, nickname, signature                      string
				gender                                            int
				commentCount, diggCount, collectCount, shareCount int
				playCount, admireCount                            int
			)

			// 扫描行数据
			err := rows.Scan(
				&videoID, &awemeID, &desc, &createTime, &authorUserID, &duration,
				&userUID, &nickname, &gender, &signature,
				&commentCount, &diggCount, &collectCount, &shareCount,
				&playCount, &admireCount,
			)
			if err != nil {
				return nil, fmt.Errorf("解析私有视频数据失败: %v", err)
			}

			// 获取用户头像
			avatarQuery := `
				SELECT uri_path, url_path, cover_type FROM cover_urls 
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
				SELECT uri, url, width, height, data_size, file_hash FROM video_play_addresses
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
				FROM music m
				JOIN videos vm ON m.id = vm.music_id
				WHERE vm.id = $1
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
				AwemeID:         awemeID,
				Desc:            desc,
				CreateTime:      createTime,
				ShareURL:        fmt.Sprintf("https://example.com/share/%s", awemeID),
				Duration:        duration,
				AuthorUserID:    authorUserID,
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
					Height:         height,
					Width:          width,
					Ratio:          "540p",
					UseStaticCover: false,
					Duration:       duration,
				},
				Statistics: Statistics{
					AdmireCount:  admireCount,
					CommentCount: commentCount,
					DiggCount:    diggCount,
					CollectCount: collectCount,
					PlayCount:    playCount,
					ShareCount:   shareCount,
				},
				Status: StatusInfo{
					ListenVideoStatus: 0,
					IsDelete:          false,
					AllowShare:        true,
					IsProhibited:      false,
					InReviewing:       false,
					PartSee:           0,
					PrivateStatus:     0,
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
					UID:       userUID,
					Nickname:  nickname,
					Gender:    gender,
					Signature: signature,
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
			return nil, fmt.Errorf("查询私有视频数据时发生错误: %v", err)
		}

		return videos, nil
	}

	// 如果没有使用数据库，返回错误
	return nil, fmt.Errorf("数据库未初始化")
}

// 获取喜欢的视频列表
func GetLikedVideosFromDB(offset, limit int) ([]Video, error) {
	if config.DB != nil {
		// 从 PostgreSQL 数据库中获取喜欢的视频数据
		query := `
			SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration,
			       u.uid, u.nickname, u.gender, u.signature, 
			       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count,
			       vs.play_count, vs.collect_count
			FROM videos v
			LEFT JOIN users u ON v.author_user_id = u.uid
			LEFT JOIN video_statistics vs ON v.id = vs.video_id
			LEFT JOIN user_like_videos ulv ON v.id = ulv.video_id
			WHERE ulv.id = (SELECT id FROM users WHERE uid = 'current_user_id') -- 实际应该使用当前登录用户的ID
			ORDER BY ulv.created_at DESC
			LIMIT $1 OFFSET $2
		`

		// 执行查询
		rows, err := config.DB.Query(query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("查询喜欢的视频数据失败: %v", err)
		}
		defer rows.Close()

		// 解析结果
		var videos []Video
		for rows.Next() {
			var (
				videoID, awemeID, desc, authorUserID              string
				createTime                                        int64
				duration                                          int
				userUID, nickname, signature                      string
				gender                                            int
				commentCount, diggCount, collectCount, shareCount int
				playCount, admireCount                            int
			)

			// 扫描行数据
			err := rows.Scan(
				&videoID, &awemeID, &desc, &createTime, &authorUserID, &duration,
				&userUID, &nickname, &gender, &signature,
				&commentCount, &diggCount, &collectCount, &shareCount,
				&playCount, &admireCount,
			)
			if err != nil {
				return nil, fmt.Errorf("解析喜欢的视频数据失败: %v", err)
			}

			// 获取用户头像
			avatarQuery := `
				SELECT uri_path, url_path, cover_type FROM cover_urls 
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
				SELECT uri, url, width, height, data_size, file_hash FROM video_play_addresses
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
				FROM music m
				JOIN videos vm ON m.id = vm.music_id
				WHERE vm.id = $1
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
				AwemeID:         awemeID,
				Desc:            desc,
				CreateTime:      createTime,
				ShareURL:        fmt.Sprintf("https://example.com/share/%s", awemeID),
				Duration:        duration,
				AuthorUserID:    authorUserID,
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
					Height:         height,
					Width:          width,
					Ratio:          "540p",
					UseStaticCover: false,
					Duration:       duration,
				},
				Statistics: Statistics{
					AdmireCount:  admireCount,
					CommentCount: commentCount,
					DiggCount:    diggCount,
					CollectCount: collectCount,
					PlayCount:    playCount,
					ShareCount:   shareCount,
				},
				Status: StatusInfo{
					ListenVideoStatus: 0,
					IsDelete:          false,
					AllowShare:        true,
					IsProhibited:      false,
					InReviewing:       false,
					PartSee:           0,
					PrivateStatus:     0,
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
					UID:       userUID,
					Nickname:  nickname,
					Gender:    gender,
					Signature: signature,
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
			return nil, fmt.Errorf("查询喜欢的视频数据时发生错误: %v", err)
		}

		return videos, nil
	}

	// 如果没有使用数据库，返回错误
	return nil, fmt.Errorf("数据库未初始化")
}

// 获取我的视频列表
func GetMyVideosFromDB(offset, limit int) ([]Video, error) {
	if config.DB != nil {
		// 从 PostgreSQL 数据库中获取我的视频数据
		query := `
			SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.video_type, 
			       u.uid, u.nickname, u.gender, u.signature, 
			       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count
			FROM videos v
			LEFT JOIN users u ON v.author_user_id = u.uid
			LEFT JOIN video_statistics vs ON v.id = vs.video_id
			WHERE v.author_user_id = 'current_user_id' -- 实际应该使用当前登录用户的ID
			ORDER BY v.create_time DESC
			LIMIT $1 OFFSET $2
		`

		// 执行查询
		rows, err := config.DB.Query(query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("查询我的视频数据失败: %v", err)
		}
		defer rows.Close()

		// 解析结果
		var videos []Video
		for rows.Next() {
			var (
				videoID, awemeID, desc, authorUserID              string
				createTime                                        int64
				duration                                          int
				userUID, nickname, signature, gender              string
				commentCount, diggCount, collectCount, shareCount int
			)

			// 扫描行数据
			err := rows.Scan(
				&videoID, &awemeID, &desc, &createTime, &authorUserID, &duration,
				&userUID, &nickname, &gender, &signature,
				&commentCount, &diggCount, &collectCount, &shareCount,
			)
			if err != nil {
				return nil, fmt.Errorf("解析我的视频数据失败: %v", err)
			}

			// 获取用户头像
			avatarQuery := `
				SELECT url_path, cover_type FROM cover_urls 
				WHERE user_id = (SELECT id FROM users WHERE uid = $1) 
				AND (cover_type = 'avatar_168x168' OR cover_type = 'avatar_300x300')
			`
			avatarRows, err := config.DB.Query(avatarQuery, userUID)
			if err != nil {
				return nil, fmt.Errorf("查询用户头像失败: %v", err)
			}
			defer avatarRows.Close()

			var avatar168URL, avatar300URL string
			for avatarRows.Next() {
				var url, avatarType string
				err := avatarRows.Scan(&url, &avatarType)
				if err != nil {
					return nil, fmt.Errorf("解析用户头像失败: %v", err)
				}

				if avatarType == "avatar_168x168" {
					avatar168URL = url
				} else if avatarType == "avatar_300x300" {
					avatar300URL = url
				}
			}

			// 构建视频对象
			video := Video{
				AwemeID:         awemeID,
				Desc:            desc,
				CreateTime:      createTime,
				ShareURL:        fmt.Sprintf("https://example.com/share/%s", awemeID),
				Duration:        duration,
				AuthorUserID:    authorUserID,
				PreventDownload: false,
				VideoInfo: VideoInfo{
					PlayAddr: PlayAddr{
						URI:      "",
						URLList:  []string{fmt.Sprintf("https://example.com/video/%s.mp4", awemeID)},
						Width:    0,
						Height:   0,
						URLKey:   "",
						DataSize: 0,
						FileHash: "",
						FileCS:   "",
					},
					Cover: Cover{
						URI:     "",
						URLList: []string{""},
						Width:   0,
						Height:  0,
					},
					Height:         0,
					Width:          0,
					Ratio:          "540p",
					UseStaticCover: false,
					Duration:       duration,
				},
				Statistics: Statistics{
					AdmireCount:  0,
					CommentCount: commentCount,
					DiggCount:    diggCount,
					CollectCount: collectCount,
					PlayCount:    0,
					ShareCount:   shareCount,
				},
				Status: StatusInfo{
					ListenVideoStatus: 0,
					IsDelete:          false,
					AllowShare:        true,
					IsProhibited:      false,
					InReviewing:       false,
					PartSee:           0,
					PrivateStatus:     0,
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
					UID:       userUID,
					Nickname:  nickname,
					Gender:    genderToInt(gender),
					Signature: signature,
					Avatar168x168: Avatar{
						URI:     "",
						URLList: []string{avatar168URL},
						Width:   168,
						Height:  168,
					},
					Avatar300x300: Avatar{
						URI:     "",
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
			return nil, fmt.Errorf("查询我的视频数据时发生错误: %v", err)
		}

		return videos, nil
	}

	// 如果没有使用数据库，返回错误
	return nil, fmt.Errorf("数据库未初始化")
}

// 获取历史视频列表
func GetHistoryVideosFromDB(offset, limit int) ([]Video, error) {
	if config.DB != nil {
		// 从 PostgreSQL 数据库中获取历史视频数据
		query := `
			SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration, 
			       u.uid, u.nickname, u.gender, u.signature, 
			       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count, vs.play_count, vs.collect_count
			FROM videos v
			LEFT JOIN users u ON v.author_user_id = u.uid
			LEFT JOIN video_statistics vs ON v.id = vs.video_id
			LEFT JOIN user_history_videos uhv ON v.id = uhv.video_id
			WHERE uhv.id = (SELECT id FROM users WHERE uid = 'current_user_id') -- 实际应该使用当前登录用户的ID
			ORDER BY uhv.view_time DESC
			LIMIT $1 OFFSET $2
		`

		// 执行查询
		rows, err := config.DB.Query(query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("查询历史视频数据失败: %v", err)
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
				return nil, fmt.Errorf("解析历史视频数据失败: %v", err)
			}

			// 获取用户头像
			avatarQuery := `
				SELECT uri_path, url_path, cover_type FROM cover_urls 
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
				SELECT uri, url, width, height, data_size, file_hash FROM video_play_addresses
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
				FROM music m
				JOIN videos vm ON m.id = vm.music_id
				WHERE vm.id = $1
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
				AwemeID:         awemeID,
				Desc:            desc,
				CreateTime:      createTime,
				ShareURL:        fmt.Sprintf("https://example.com/share/%s", awemeID),
				Duration:        duration,
				AuthorUserID:    authorUserID,
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
					Height:         height,
					Width:          width,
					Ratio:          "540p",
					UseStaticCover: false,
					Duration:       duration,
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
					IsDelete:          false,
					AllowShare:        true,
					IsProhibited:      false,
					InReviewing:       false,
					PartSee:           0,
					PrivateStatus:     0,
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
					UID:       userUID,
					Nickname:  nickname,
					Gender:    genderToInt(gender),
					Signature: signature,
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
			return nil, fmt.Errorf("查询历史视频数据时发生错误: %v", err)
		}

		return videos, nil
	}

	// 如果没有使用数据库，返回错误
	return nil, fmt.Errorf("数据库未初始化")
}

// 获取视频总数
func GetVideoCountFromDB(videoType string) (int, error) {
	if config.DB != nil {
		// 从 PostgreSQL 数据库中获取视频总数
		query := `
			SELECT COUNT(*) FROM videos
			WHERE video_type = $1
		`

		// 执行查询
		var count int
		err := config.DB.QueryRow(query, videoType).Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("查询视频总数失败: %v", err)
		}

		return count, nil
	}

	// 如果没有使用数据库，返回错误
	return 0, fmt.Errorf("数据库未初始化")
}

// 从数据库获取视频统计信息
func GetVideoStatistics(videoID string) (DBVideoStatistics, error) {
	if config.DB == nil {
		return DBVideoStatistics{}, fmt.Errorf("数据库未初始化")
	}

	query := `
		SELECT comment_count, digg_count, collect_count, play_count, share_count
		FROM video_statistics
		WHERE video_id = (SELECT id FROM videos WHERE aweme_id = $1)
	`
	var stats DBVideoStatistics
	err := config.DB.QueryRow(query, videoID).Scan(
		&stats.CommentCount,
		&stats.DiggCount,
		&stats.CollectCount,
		&stats.PlayCount,
		&stats.ShareCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return DBVideoStatistics{}, fmt.Errorf("视频统计信息不存在")
		}
		return DBVideoStatistics{}, err
	}

	return stats, nil
}
