package model

import (
	"database/sql"
	"fmt"
	"klik/server/config"
	"log"
)

// GetRecommendPostsFromDB 从数据库获取推荐帖子
func GetRecommendPostsFromDB(offset, limit int) ([]Post, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 查询帖子数据
	query := `
		SELECT p.id, p.post_id, p.text, p.post_id, p.create_time, p.digg_count, p.comment_count, p.share_count
		FROM posts p
		ORDER BY p.create_time DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询帖子数据失败: %v", err)
	}
	defer rows.Close()

	// 处理查询结果
	var posts []Post
	var dbPosts []DBPost

	for rows.Next() {
		var post DBPost
		err := rows.Scan(
			&post.ID,
			&post.PostID,
			&post.Text,
			&post.AuthorID,
			&post.CreateTime,
			&post.DiggCount,
			&post.CommentCount,
			&post.ShareCount,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描帖子数据失败: %v", err)
		}
		dbPosts = append(dbPosts, post)
	}

	// 获取每个帖子的作者信息和图片
	for i, post := range dbPosts {
		// 获取作者信息
		author, err := getUserByIDFromDB(post.AuthorID)
		if err != nil {
			log.Printf("获取作者信息失败: %v", err)
			continue
		}
		dbPosts[i].Author = author

		// 获取帖子图片
		images, err := getPostImagesFromDB(post.ID)
		if err != nil {
			log.Printf("获取帖子图片失败: %v", err)
			continue
		}
		dbPosts[i].Images = images

		// 转换为Post模型
		posts = append(posts, Post{
			ID:           post.PostID,
			Text:         post.Text,
			Author:       post.Author,
			CreateTime:   post.CreateTime,
			DiggCount:    post.DiggCount,
			CommentCount: post.CommentCount,
			ShareCount:   post.ShareCount,
			Images:       post.Images,
		})
	}

	return posts, nil
}

// GetPostCountFromDB 获取帖子总数
func GetPostCountFromDB() (int, error) {
	if config.DB == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}

	// 查询帖子总数
	var count int
	query := "SELECT COUNT(*) FROM posts"
	err := config.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询帖子总数失败: %v", err)
	}

	return count, nil
}

// getUserByIDFromDB 从数据库获取用户信息（通过ID）
func getUserByIDFromDB(id int) (User, error) {
	if config.DB == nil {
		return User{}, fmt.Errorf("数据库未初始化")
	}

	// 查询用户数据
	query := `
		SELECT id, uid, nickname, gender, signature, 
		       avatar_168x168_uri, avatar_168x168_url, avatar_300x300_uri, avatar_300x300_url
		FROM users
		WHERE id = $1
	`
	var dbUser DBUser
	err := config.DB.QueryRow(query, id).Scan(
		&dbUser.ID,
		&dbUser.UID,
		&dbUser.Nickname,
		&dbUser.Gender,
		&dbUser.Signature,
		&dbUser.Avatar168x168URI,
		&dbUser.Avatar168x168URL,
		&dbUser.Avatar300x300URI,
		&dbUser.Avatar300x300URL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("用户不存在")
		}
		return User{}, fmt.Errorf("查询用户数据失败: %v", err)
	}

	// 构建User对象
	user := User{
		UID:       dbUser.UID,
		Nickname:  dbUser.Nickname,
		Gender:    dbUser.Gender,
		Signature: dbUser.Signature,
		Avatar168x168: AvatarInfo{
			URLList: []string{dbUser.Avatar168x168URL},
		},
		Avatar300x300: AvatarInfo{
			URLList: []string{dbUser.Avatar300x300URL},
		},
	}

	return user, nil
}

// getPostImagesFromDB 从数据库获取帖子图片
func getPostImagesFromDB(postID int) ([]string, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 查询帖子图片
	query := `
		SELECT image_url
		FROM post_images
		WHERE post_id = $1
		ORDER BY id
	`
	rows, err := config.DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("查询帖子图片失败: %v", err)
	}
	defer rows.Close()

	// 处理查询结果
	var images []string
	for rows.Next() {
		var imageURL string
		err := rows.Scan(&imageURL)
		if err != nil {
			return nil, fmt.Errorf("扫描帖子图片数据失败: %v", err)
		}
		images = append(images, imageURL)
	}

	return images, nil
}

// GetGoodsFromDB 从数据库获取商品
func GetGoodsFromDB(offset, limit int) ([]Good, error) {
	if config.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 查询商品数据
	query := `
		SELECT id, good_id, title, description, price, image, sale_count
		FROM goods
		ORDER BY sale_count DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询商品数据失败: %v", err)
	}
	defer rows.Close()

	// 处理查询结果
	var goods []Good
	for rows.Next() {
		var dbGood DBGood
		err := rows.Scan(
			&dbGood.ID,
			&dbGood.GoodID,
			&dbGood.Title,
			&dbGood.Description,
			&dbGood.Price,
			&dbGood.Image,
			&dbGood.SaleCount,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描商品数据失败: %v", err)
		}

		// 转换为Good模型
		goods = append(goods, Good{
			ID:          dbGood.GoodID,
			Title:       dbGood.Title,
			Description: dbGood.Description,
			Price:       dbGood.Price,
			Image:       dbGood.Image,
			SaleCount:   dbGood.SaleCount,
		})
	}

	return goods, nil
}

// GetGoodCountFromDB 获取商品总数
func GetGoodCountFromDB() (int, error) {
	if config.DB == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}

	// 查询商品总数
	var count int
	query := "SELECT COUNT(*) FROM goods"
	err := config.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询商品总数失败: %v", err)
	}

	return count, nil
}
