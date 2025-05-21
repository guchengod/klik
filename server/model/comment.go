package model

import (
	"database/sql"
	"fmt"
	"klik/server/config"
)

// genderToInt 将性别字符串转换为整数
func genderToInt(gender string) int {
	switch gender {
	case "男", "male", "M":
		return 1
	case "女", "female", "F":
		return 2
	default:
		return 0
	}
}

// 从数据库获取视频评论
func GetVideoCommentsFromPostgres(videoID string) ([]Comment, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `
		SELECT c.id, c.text, c.create_time, c.digg_count, c.parent_id
		FROM comments c
		JOIN videos v ON c.video_id = v.id
		WHERE v.aweme_id = $1
		ORDER BY c.create_time DESC
	`
	rows, err := config.DB.Query(query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var userID string
		err := rows.Scan(
			&comment.ID,
			&comment.Text,
			&comment.CreateTime,
			&comment.DiggCount,
			&userID,
		)
		if err != nil {
			return nil, err
		}

		// 查询用户信息
		comment.User, err = GetUserByID(userID)
		if err != nil {
			// 如果用户不存在，使用默认用户
			comment.User = User{
				UID:      userID,
				Nickname: "未知用户",
			}
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// GetLongRecommendVideosFromPostgres 从数据库获取长视频推荐列表
func GetLongRecommendVideosFromPostgres(start, pageSize int) ([]Video, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 查询视频列表
	query := `
		SELECT v.id, v.aweme_id, v.video_desc, v.create_time, v.author_user_id, v.duration,
		       u.uid, u.nickname, u.gender, u.signature,
		       vs.comment_count, vs.digg_count, vs.collect_count, vs.share_count, vs.play_count, vs.admire_count
		FROM videos v
		LEFT JOIN users u ON v.author_user_id = u.uid
		LEFT JOIN video_statistics vs ON v.id = vs.video_id
		WHERE v.duration > 60
		ORDER BY v.create_time DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := config.DB.Query(query, pageSize, start)
	if err != nil {
		return nil, err
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

		err := rows.Scan(
			&videoID,
			&awemeID,
			&desc,
			&createTime,
			&authorUserID,
			&duration,
			&userUID,
			&nickname,
			&gender,
			&signature,
			&commentCount,
			&diggCount,
			&collectCount,
			&shareCount,
			&playCount,
			&admireCount,
		)
		if err != nil {
			return nil, err
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

		// 获取视频播放地址
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

	return videos, nil
}
