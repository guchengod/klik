-- 导入数据示例脚本
-- 注意：实际导入应该使用程序来处理复杂的JSON数据

-- 连接到数据库
\c klik;

-- 创建导入函数
CREATE OR REPLACE FUNCTION import_posts6_data() RETURNS void AS $$
DECLARE
    current_user_id INTEGER;
    current_music_id BIGINT;
    current_video_id INTEGER;
BEGIN
    -- 清空现有数据（如果需要）
    -- TRUNCATE users, cover_urls, music, videos, video_play_addresses, video_covers, video_statistics, video_status, comments CASCADE;
    
    -- 这里应该是实际导入逻辑
    -- 由于JSON数据结构复杂，建议使用编程语言（如Go、Python等）来处理导入
    -- 下面是一个简化的示例，展示如何导入一些基本数据
    
    -- 插入示例用户
    INSERT INTO users (uid, nickname, gender, signature) 
    VALUES ('59054327754', '我是香秀🐂🍺', 2, '合作：X229896（备注品牌 ）')
    RETURNING id INTO current_user_id;
    
    -- 插入用户头像
    INSERT INTO cover_urls (user_id, uri, url, type)
    VALUES 
    (current_user_id, 'aweme-avatar/mosaic-legacy_20b7700050147c01968f3', 'https://p3-pc.douyinpic.com/img/aweme-avatar/mosaic-legacy_20b7700050147c01968f3~c5_168x168.jpeg?from=2956013662', 'avatar_168x168'),
    (current_user_id, 'aweme-avatar/mosaic-legacy_20b7700050147c01968f3', 'https://p3-pc.douyinpic.com/img/aweme-avatar/mosaic-legacy_20b7700050147c01968f3~c5_300x300.jpeg?from=2956013662', 'avatar_300x300'),
    (current_user_id, 'douyin-user-file/4eec4c18569133d5990381a62ba49327', 'fmO_JqQD-ukKguwbdyoiL.png', 'cover'),
    (current_user_id, 'douyin-user-file/4eec4c18569133d5990381a62ba49327', 'OvKvfthk8TXKeVpwEkQNq.png', 'white_cover');
    
    -- 插入示例音乐
    INSERT INTO music (id, title, author, duration, owner_nickname, is_original)
    VALUES (7123453673090321000, '禁盗用', 'LoveW_', 17, 'LoveW_', false)
    RETURNING id INTO current_music_id;
    
    -- 插入示例视频
    INSERT INTO videos (aweme_id, desc, create_time, music_id, author_user_id, duration, type, is_top)
    VALUES ('7260749400622894336', '你说爱像云 要自在漂浮才美丽', 1690524964, current_music_id, '59054327754', 13560, 'recommend-video', true)
    RETURNING id INTO current_video_id;
    
    -- 插入视频播放地址
    INSERT INTO video_play_addresses (video_id, uri, url, width, height, data_size, file_hash)
    VALUES (current_video_id, 'v0d00fg10000cj1lq4jc77u0ng6s1gt0', 'https://www.douyin.com/aweme/v1/play/?video_id=v0d00fg10000cj1lq4jc77u0ng6s1gt0&line=0&file_id=bed51c00899b458cbc5d8280147c22a1&sign=7749aec7bd62a3760065f60e40fc1867&is_play_url=1&source=PackSourceEnum_PUBLISH', 1080, 1920, 3480280, '7749aec7bd62a3760065f60e40fc1867');
    
    -- 插入视频封面
    INSERT INTO video_covers (video_id, uri, url, width, height)
    VALUES (current_video_id, 'tos-cn-i-0813/oYVDeaFZyENAAAAKXCYfxD6hI4zADNAURgtySl', 'jwWCPZVTIA4IKM-8WipLF.png', 720, 720);
    
    -- 插入视频统计
    INSERT INTO video_statistics (video_id, comment_count, digg_count, collect_count, play_count, share_count)
    VALUES (current_video_id, 21582, 1246636, 64460, 0, 172803);
    
    -- 插入视频状态
    INSERT INTO video_status (video_id, is_delete, allow_share, is_prohibited, in_reviewing, private_status)
    VALUES (current_video_id, false, true, false, false, 0);
    
    -- 更多数据导入...
    -- 实际应用中，应该使用程序读取JSON文件并导入所有数据
END;
$$ LANGUAGE plpgsql;

-- 执行导入函数
SELECT import_posts6_data();
