-- 导入视频数据脚本
-- 连接到数据库
\c douyin;

BEGIN;

-- ==================== 导入第一个视频记录 ====================
-- 插入用户
INSERT INTO users (uid, nickname, gender, signature, ip_location, follower_count, following_count, total_favorited, 
                  aweme_count, birthday_hide_level, can_show_group_card, commerce_user_level)
VALUES ('59054327754', '我是香秀🐂🍺', 2, '合作：X229896（备注品牌 ）', 'IP属地：天津', 7078268, 88, 202309485, 
        296, 1, 1, 0)
ON CONFLICT (uid) DO UPDATE SET 
    nickname = EXCLUDED.nickname,
    gender = EXCLUDED.gender,
    signature = EXCLUDED.signature,
    ip_location = EXCLUDED.ip_location,
    follower_count = EXCLUDED.follower_count,
    following_count = EXCLUDED.following_count;

-- 插入用户头像
INSERT INTO cover_urls (user_id, uri_path, url_path, cover_type)
SELECT id, 'aweme-avatar/mosaic-legacy_20b7700050147c01968f3', 
       'https://p3-pc.douyinpic.com/img/aweme-avatar/mosaic-legacy_20b7700050147c01968f3~c5_168x168.jpeg?from=2956013662', 
       'avatar_168x168'
FROM users WHERE uid = '59054327754'
ON CONFLICT DO NOTHING;

INSERT INTO cover_urls (user_id, uri_path, url_path, cover_type)
SELECT id, 'aweme-avatar/mosaic-legacy_20b7700050147c01968f3', 
       'https://p3-pc.douyinpic.com/img/aweme-avatar/mosaic-legacy_20b7700050147c01968f3~c5_300x300.jpeg?from=2956013662', 
       'avatar_300x300'
FROM users WHERE uid = '59054327754'
ON CONFLICT DO NOTHING;

-- 插入商业信息
INSERT INTO commerce_user_info (user_id, has_ads_entry, show_star_atlas_cooperation, star_atlas)
SELECT id, true, true, 1
FROM users WHERE uid = '59054327754'
ON CONFLICT DO NOTHING;

-- 插入音乐
INSERT INTO music (id, id_str, title, author, duration, owner_id, owner_nickname, is_original, play_url)
VALUES (7123453673090321000, '7123453673090321159', '禁盗用', 'LoveW_', 17, '1711265034548715', 'LoveW_', false, 
        'https://sf5-hl-cdn-tos.douyinstatic.com/obj/ies-music/7123453672335690532.mp3')
ON CONFLICT (id) DO UPDATE SET 
    title = EXCLUDED.title,
    author = EXCLUDED.author,
    duration = EXCLUDED.duration,
    owner_nickname = EXCLUDED.owner_nickname;

-- 插入视频
INSERT INTO videos (aweme_id, video_desc, create_time, author_user_id, duration, music_id, share_url, is_top, prevent_download)
VALUES ('7260749400622894336', '你说爱像云 要自在漂浮才美丽', 1690524964, '59054327754', 13560, 7123453673090321000,
        'https://www.iesdouyin.com/share/video/7260749400622894336/?region=CN&mid=7123453673090321159', 
        true, false)
ON CONFLICT (aweme_id) DO UPDATE SET 
    video_desc = EXCLUDED.video_desc,
    create_time = EXCLUDED.create_time,
    author_user_id = EXCLUDED.author_user_id,
    duration = EXCLUDED.duration;

-- 获取插入的视频ID
DO $$
DECLARE
    video_id INTEGER;
BEGIN
    SELECT id INTO video_id FROM videos WHERE aweme_id = '7260749400622894336';

    -- 插入视频播放地址
    INSERT INTO video_play_addresses (video_id, uri, url, width, height, data_size, file_hash)
    VALUES (video_id, 'v0d00fg10000cj1lq4jc77u0ng6s1gt0', 
            'https://www.douyin.com/aweme/v1/play/?video_id=v0d00fg10000cj1lq4jc77u0ng6s1gt0&line=0&file_id=bed51c00899b458cbc5d8280147c22a1&sign=7749aec7bd62a3760065f60e40fc1867&is_play_url=1&source=PackSourceEnum_PUBLISH',
            1080, 1920, 3480280, '7749aec7bd62a3760065f60e40fc1867')
    ON CONFLICT DO NOTHING;

    -- 插入视频封面
    INSERT INTO video_covers (video_id, uri, url, width, height)
    VALUES (video_id, 'tos-cn-i-0813/oYVDeaFZyENAAAAKXCYfxD6hI4zADNAURgtySl', 'jwWCPZVTIA4IKM-8WipLF.png', 720, 720)
    ON CONFLICT DO NOTHING;

    -- 插入视频统计信息
    INSERT INTO video_statistics (video_id, comment_count, digg_count, collect_count, play_count, share_count)
    VALUES (video_id, 21582, 1246636, 64460, 0, 172803)
    ON CONFLICT DO NOTHING;

    -- 插入视频状态
    INSERT INTO video_status (video_id, is_delete, allow_share, is_prohibited, in_reviewing, private_status)
    VALUES (video_id, false, true, false, false, 0)
    ON CONFLICT DO NOTHING;
END
$$;

-- ==================== 导入第二个视频记录 ====================
-- 插入用户 (假设这是另一个用户)
INSERT INTO users (uid, nickname, gender, signature)
VALUES ('62839305427', '普通用户', 0, '这是一个普通用户')
ON CONFLICT (uid) DO UPDATE SET 
    nickname = EXCLUDED.nickname,
    gender = EXCLUDED.gender,
    signature = EXCLUDED.signature;

-- 插入音乐
INSERT INTO music (id, id_str, title, author, duration, is_original)
VALUES (6452110567468436000, '6452110567468436238', 'GQ', 'Lola Coca', 53, false)
ON CONFLICT (id) DO UPDATE SET 
    title = EXCLUDED.title,
    author = EXCLUDED.author,
    duration = EXCLUDED.duration;

-- 插入视频
INSERT INTO videos (aweme_id, video_desc, create_time, author_user_id, duration, music_id, share_url, is_top, prevent_download)
VALUES ('6686589698707590411', '门有点矮哟～', 1556887936, '62839305427', 7133, 6452110567468436000,
        'https://www.iesdouyin.com/share/video/6686589698707590411/?region=CN&mid=6452110567468436238', 
        true, false)
ON CONFLICT (aweme_id) DO UPDATE SET 
    video_desc = EXCLUDED.video_desc,
    create_time = EXCLUDED.create_time,
    author_user_id = EXCLUDED.author_user_id,
    duration = EXCLUDED.duration;

-- 获取插入的视频ID
DO $$
DECLARE
    video_id INTEGER;
BEGIN
    SELECT id INTO video_id FROM videos WHERE aweme_id = '6686589698707590411';

    -- 插入视频播放地址
    INSERT INTO video_play_addresses (video_id, uri, url, width, height, data_size, file_hash)
    VALUES (video_id, 'v0200f2c0000bj63fuv3cp5a1sbmujc0', 
            'https://www.douyin.com/aweme/v1/play/?video_id=v0200f2c0000bj63fuv3cp5a1sbmujc0&line=0&file_id=fad24ab3e1ab4efa90440244d5268bd9&sign=f33a08757b00e73f6a75ab6a3c5d751b&is_play_url=1&source=PackSourceEnum_PUBLISH',
            720, 1280, 1578648, 'f33a08757b00e73f6a75ab6a3c5d751b')
    ON CONFLICT DO NOTHING;

    -- 插入视频封面
    INSERT INTO video_covers (video_id, uri, url, width, height)
    VALUES (video_id, '2071800047c9d7f961529', '_T0vQPZKXrNC6ulECmMes.png', 720, 720)
    ON CONFLICT DO NOTHING;

    -- 插入视频统计信息
    INSERT INTO video_statistics (video_id, comment_count, digg_count, collect_count, play_count, share_count)
    VALUES (video_id, 19109, 865701, 11578, 0, 44504)
    ON CONFLICT DO NOTHING;

    -- 插入视频状态
    INSERT INTO video_status (video_id, is_delete, allow_share, is_prohibited, in_reviewing, private_status)
    VALUES (video_id, false, true, false, false, 0)
    ON CONFLICT DO NOTHING;
END
$$;

-- 创建一些示例用户的收藏和点赞记录
DO $$
DECLARE
    current_user_id VARCHAR(50) := '59054327754';
    video_id1 INTEGER;
    video_id2 INTEGER;
BEGIN
    SELECT id INTO video_id1 FROM videos WHERE aweme_id = '7260749400622894336';
    SELECT id INTO video_id2 FROM videos WHERE aweme_id = '6686589698707590411';
    
    -- 用户收藏视频
    INSERT INTO user_collect_videos (commenter_id, video_id)
    VALUES (current_user_id, video_id2)
    ON CONFLICT DO NOTHING;
    
    -- 用户喜欢视频
    INSERT INTO user_like_videos (commenter_id, video_id)
    VALUES (current_user_id, video_id1)
    ON CONFLICT DO NOTHING;
    
    -- 用户历史记录
    INSERT INTO user_history_videos (commenter_id, video_id, view_time)
    VALUES (current_user_id, video_id1, NOW() - INTERVAL '1 day'),
           (current_user_id, video_id2, NOW() - INTERVAL '2 days')
    ON CONFLICT DO NOTHING;
END
$$;

-- 添加示例评论
DO $$
DECLARE
    video_id1 INTEGER;
    user_id VARCHAR(50) := '59054327754';
BEGIN
    SELECT id INTO video_id1 FROM videos WHERE aweme_id = '7260749400622894336';
    
    -- 添加视频评论
    INSERT INTO comments (comment_id, video_id, commenter_id, content, create_time, digg_count)
    VALUES ('comment_1', video_id1, user_id, '这个视频真不错！', extract(epoch from now())::bigint, 123),
           ('comment_2', video_id1, user_id, '非常喜欢这个内容', extract(epoch from now())::bigint, 45)
    ON CONFLICT DO NOTHING;
    
    -- 添加子评论
    INSERT INTO sub_comments (comment_id, parent_cmt_id, commenter_id, content, create_time, digg_count)
    VALUES ('sub_comment_1', 'comment_1', user_id, '我也这么认为！', extract(epoch from now())::bigint, 56)
    ON CONFLICT DO NOTHING;
END
$$;

-- 创建一些长视频示例
INSERT INTO videos (aweme_id, video_desc, create_time, author_user_id, duration, video_type, share_url)
VALUES ('long_video_1', '这是一个长视频示例', extract(epoch from now())::bigint, '59054327754', 300, 'long-video', 'https://example.com/share/long_video_1'),
       ('long_video_2', '另一个长视频示例', extract(epoch from now())::bigint, '62839305427', 500, 'long-video', 'https://example.com/share/long_video_2')
ON CONFLICT DO NOTHING;

-- 创建一些其他类型的视频
INSERT INTO videos (aweme_id, video_desc, create_time, author_user_id, duration, video_type, share_url)
VALUES ('private_video_1', '这是一个私人视频', extract(epoch from now())::bigint, '59054327754', 120, 'private-video', 'https://example.com/share/private_video_1'),
       ('liked_video_1', '这是一个喜欢的视频', extract(epoch from now())::bigint, '62839305427', 60, 'liked-video', 'https://example.com/share/liked_video_1')
ON CONFLICT DO NOTHING;

-- 为所有新增视频添加统计信息
DO $$
DECLARE
    v_id INTEGER;
    v_record RECORD;
BEGIN
    FOR v_record IN SELECT id FROM videos WHERE id NOT IN (SELECT video_id FROM video_statistics) LOOP
        v_id := v_record.id;
        
        -- 插入视频统计信息
        INSERT INTO video_statistics (video_id, comment_count, digg_count, collect_count, play_count, share_count)
        VALUES (v_id, floor(random() * 1000)::int, floor(random() * 10000)::int, 
                floor(random() * 500)::int, floor(random() * 50000)::int, floor(random() * 200)::int)
        ON CONFLICT DO NOTHING;
        
        -- 插入视频状态
        INSERT INTO video_status (video_id, is_delete, allow_share, is_prohibited, in_reviewing, private_status)
        VALUES (v_id, false, true, false, false, 0)
        ON CONFLICT DO NOTHING;
    END LOOP;
END
$$;

COMMIT;
