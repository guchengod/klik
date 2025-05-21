-- 创建数据库
CREATE DATABASE douyin;

-- 连接到数据库
\c douyin;

-- 创建用户表
CREATE TABLE users
(
    id                 SERIAL PRIMARY KEY,
    uid                VARCHAR(50) UNIQUE NOT NULL,
    sec_uid            VARCHAR(100),
    short_user_id      VARCHAR(50),
    user_unique_id     VARCHAR(100),
    nickname           VARCHAR(100)       NOT NULL,
    avatar_168x168_uri VARCHAR(255),
    avatar_168x168_url TEXT,
    avatar_300x300_uri VARCHAR(255),
    avatar_300x300_url TEXT,
    gender             INTEGER                  DEFAULT 0,
    signature          TEXT,
    ip_location        VARCHAR(100),
    province           VARCHAR(100),
    city               VARCHAR(100),
    country            VARCHAR(100),
    district           VARCHAR(100),
    birthday_hide_level INTEGER                 DEFAULT 0,
    can_show_group_card INTEGER                 DEFAULT 1,
    commerce_user_level INTEGER                 DEFAULT 0,
    cover_colour       VARCHAR(20),
    favoriting_count   INTEGER                  DEFAULT 0,
    follow_status      INTEGER                  DEFAULT 0,
    follower_count     INTEGER                  DEFAULT 0,
    follower_request_status INTEGER              DEFAULT 0,
    follower_status    INTEGER                  DEFAULT 0,
    following_count    INTEGER                  DEFAULT 0,
    forward_count      INTEGER                  DEFAULT 0,
    max_follower_count INTEGER                  DEFAULT 0,
    mplatform_followers_count INTEGER           DEFAULT 0,
    public_collects_count INTEGER               DEFAULT 0,
    total_favorited    INTEGER                  DEFAULT 0,
    aweme_count        INTEGER                  DEFAULT 0,
    unique_id          VARCHAR(100),
    short_id           VARCHAR(50),
    user_age           INTEGER                  DEFAULT -1,
    created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建封面URL表
CREATE TABLE cover_urls
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    uri_path    VARCHAR(255),                         -- 将 'uri' 改为 'uri_path'，'uri' 可能是保留字
    url_path    TEXT,                                 -- 将 'url' 改为 'url_path'，'url' 可能是保留字
    width      INTEGER                  DEFAULT 720,
    height     INTEGER                  DEFAULT 720,
    cover_type VARCHAR(20) NOT NULL, -- 'cover' 或 'white_cover'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户商业信息表
CREATE TABLE commerce_user_info
(
    id                      SERIAL PRIMARY KEY,
    user_id                 INTEGER REFERENCES users (id) ON DELETE CASCADE,
    has_ads_entry           BOOLEAN                  DEFAULT FALSE,
    show_star_atlas_cooperation BOOLEAN             DEFAULT FALSE,
    star_atlas              INTEGER                  DEFAULT 0,
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户卡片条目表
CREATE TABLE card_entries
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER REFERENCES users (id) ON DELETE CASCADE,
    goto_url    TEXT,
    sub_title   VARCHAR(100),
    title       VARCHAR(100),
    entry_type  INTEGER                  DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户分享信息表
CREATE TABLE share_info
(
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER REFERENCES users (id) ON DELETE CASCADE,
    persist_flag   INTEGER                  DEFAULT 1,  -- 将 'bool_persist' 改为 'persist_flag'，'bool' 可能是保留字
    share_desc     TEXT,
    share_title    TEXT,
    share_url      TEXT,
    share_weibo_desc TEXT,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建音乐表
CREATE TABLE music
(
    id             BIGINT PRIMARY KEY,
    id_str         VARCHAR(50),
    title          VARCHAR(255) NOT NULL,
    author         VARCHAR(255),
    album          VARCHAR(255),
    cover_uri      VARCHAR(255),
    cover_url      TEXT,
    play_url       TEXT,
    duration       INTEGER                  DEFAULT 0,
    owner_id       VARCHAR(50),
    owner_nickname VARCHAR(100),
    is_original    BOOLEAN                  DEFAULT FALSE,
    source_platform INTEGER                 DEFAULT 0,
    is_restricted  BOOLEAN                  DEFAULT FALSE,
    is_video_self_see BOOLEAN              DEFAULT FALSE,
    prevent_download BOOLEAN               DEFAULT FALSE,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建视频表
CREATE TABLE videos
(
    id               SERIAL PRIMARY KEY,
    aweme_id         VARCHAR(50) UNIQUE NOT NULL,
    video_desc       TEXT,
    create_time      BIGINT,
    music_id         BIGINT REFERENCES music (id),
    author_user_id   VARCHAR(50) REFERENCES users (uid),
    duration         INTEGER                  DEFAULT 0,
    video_type       VARCHAR(50)              DEFAULT 'recommend-video',
    share_url        TEXT,
    is_top           BOOLEAN                  DEFAULT FALSE,
    prevent_download BOOLEAN                  DEFAULT FALSE,
    is_ads           BOOLEAN                  DEFAULT FALSE,
    is_hash_tag      BOOLEAN                  DEFAULT FALSE,
    region           VARCHAR(50),
    video_labels     TEXT,
    sort_label       VARCHAR(50),
    mark_largely_following BOOLEAN           DEFAULT FALSE,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建视频播放地址表
CREATE TABLE video_play_addresses
(
    id         SERIAL PRIMARY KEY,
    video_id   INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    uri        VARCHAR(255),
    url        TEXT,
    width      INTEGER                  DEFAULT 0,
    height     INTEGER                  DEFAULT 0,
    data_size  BIGINT                   DEFAULT 0,
    file_hash  VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建视频封面表
CREATE TABLE video_covers
(
    id         SERIAL PRIMARY KEY,
    video_id   INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    uri        VARCHAR(255),
    url        TEXT,
    width      INTEGER                  DEFAULT 0,
    height     INTEGER                  DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建视频统计表
CREATE TABLE video_statistics
(
    id            SERIAL PRIMARY KEY,
    video_id      INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    comment_count INTEGER                  DEFAULT 0,
    digg_count    INTEGER                  DEFAULT 0,
    collect_count INTEGER                  DEFAULT 0,
    play_count    INTEGER                  DEFAULT 0,
    share_count   INTEGER                  DEFAULT 0,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建视频状态表
CREATE TABLE video_status
(
    id             SERIAL PRIMARY KEY,
    video_id       INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    is_delete      BOOLEAN                  DEFAULT FALSE,
    allow_share    BOOLEAN                  DEFAULT TRUE,
    is_prohibited  BOOLEAN                  DEFAULT FALSE,
    in_reviewing   BOOLEAN                  DEFAULT FALSE,
    private_status INTEGER                  DEFAULT 0,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建帖子表
CREATE TABLE posts
(
    id             SERIAL PRIMARY KEY,
    aweme_id       VARCHAR(50) UNIQUE NOT NULL,
    post_id        VARCHAR(50) UNIQUE NOT NULL,
    description    TEXT,                       -- 将 'desc' 改为 'description'，因为 'desc' 是SQL关键字
    post_text      TEXT,
    author_user_id VARCHAR(50) REFERENCES users (uid),
    create_time    BIGINT,
    duration       INTEGER                  DEFAULT 0,
    is_top         INTEGER                  DEFAULT 0,
    admire_count   INTEGER                  DEFAULT 0,
    digg_count     INTEGER                  DEFAULT 0,
    collect_count  INTEGER                  DEFAULT 0,
    play_count     INTEGER                  DEFAULT 0,
    comment_count  INTEGER                  DEFAULT 0,
    share_count    INTEGER                  DEFAULT 0,
    prevent_download BOOLEAN               DEFAULT FALSE,
    horizontal_type INTEGER                 DEFAULT 1,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建评论表
CREATE TABLE comments
(
    id           SERIAL PRIMARY KEY,
    comment_id   VARCHAR(50) UNIQUE NOT NULL,
    aweme_id     VARCHAR(50),
    video_id     INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    post_id      INTEGER REFERENCES posts (id) ON DELETE CASCADE,   -- 修复了 NULL 位置问题
    commenter_id VARCHAR(50) REFERENCES users (uid),
    content      TEXT               NOT NULL,
    ip_location  VARCHAR(100),
    create_time  BIGINT,
    digg_count   INTEGER                  DEFAULT 0,
    user_digged  INTEGER                  DEFAULT 0,
    is_author_digged BOOLEAN             DEFAULT FALSE,
    is_hot       BOOLEAN                  DEFAULT FALSE,
    is_folded    BOOLEAN                  DEFAULT FALSE,
    user_buried  BOOLEAN                  DEFAULT FALSE,
    sub_comment_count INTEGER             DEFAULT 0,
    last_modify_ts BIGINT,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT comment_target_check CHECK (
        (video_id IS NOT NULL AND post_id IS NULL) OR
        (video_id IS NULL AND post_id IS NOT NULL)
    )
);

-- 创建帖子图片表
CREATE TABLE post_images
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    image_url  TEXT,
    width      INTEGER                  DEFAULT 0,
    height     INTEGER                  DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建帖子视频表
CREATE TABLE post_videos
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    uri        VARCHAR(255),
    url        TEXT,
    width      INTEGER                  DEFAULT 0,
    height     INTEGER                  DEFAULT 0,
    ratio      VARCHAR(20),
    data_size  BIGINT                   DEFAULT 0,
    file_hash  VARCHAR(255),
    use_static_cover BOOLEAN            DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建帖子封面表
CREATE TABLE post_covers
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    uri        VARCHAR(255),
    url        TEXT,
    width      INTEGER                  DEFAULT 0,
    height     INTEGER                  DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建帖子状态表
CREATE TABLE post_status
(
    id             SERIAL PRIMARY KEY,
    post_id        INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    listen_video_status INTEGER             DEFAULT 0,
    is_delete      BOOLEAN                  DEFAULT FALSE,
    allow_share    BOOLEAN                  DEFAULT TRUE,
    is_prohibited  BOOLEAN                  DEFAULT FALSE,
    in_reviewing   BOOLEAN                  DEFAULT FALSE,
    part_see       INTEGER                  DEFAULT 0,
    private_status INTEGER                  DEFAULT 0,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建帖子标签表
CREATE TABLE post_text_extra
(
    id             SERIAL PRIMARY KEY,
    post_id        INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    start_index    INTEGER                  DEFAULT 0,
    end_index      INTEGER                  DEFAULT 0,
    extra_type     INTEGER                  DEFAULT 0,
    hashtag_name   VARCHAR(100),
    hashtag_id     VARCHAR(50),
    is_commerce    BOOLEAN                  DEFAULT FALSE,
    caption_start  INTEGER                  DEFAULT 0,
    caption_end    INTEGER                  DEFAULT 0,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建帖子分享信息表
CREATE TABLE post_share_info
(
    id             SERIAL PRIMARY KEY,
    post_id        INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    share_url      TEXT,
    share_link_desc TEXT,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建帖子视频标签表
CREATE TABLE post_video_tags
(
    id             SERIAL PRIMARY KEY,
    post_id        INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    tag_id         INTEGER                  DEFAULT 0,
    tag_name       VARCHAR(100),
    level          INTEGER                  DEFAULT 0,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建商品表
CREATE TABLE goods
(
    id          SERIAL PRIMARY KEY,
    good_id     VARCHAR(50) UNIQUE NOT NULL,
    title       VARCHAR(255)       NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2)           DEFAULT 0,
    image       TEXT,
    sale_count  INTEGER                  DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户收藏音乐关系表
CREATE TABLE user_collect_music
(
    id           SERIAL PRIMARY KEY,
    commenter_id VARCHAR(50) REFERENCES users (uid),
    music_id     BIGINT REFERENCES music (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (commenter_id, music_id)
);

-- 创建小红书笔记表
CREATE TABLE xhs_notes
(
    id             SERIAL PRIMARY KEY,
    note_id        VARCHAR(50) UNIQUE NOT NULL,
    model_type     VARCHAR(20)              DEFAULT 'note',
    display_title  TEXT,
    note_type      VARCHAR(20)              DEFAULT 'normal',  -- 将 'type' 改为 'note_type'，因为 'type' 是常见SQL关键字
    author_user_id VARCHAR(50) REFERENCES users (uid),
    liked_count    INTEGER                  DEFAULT 0,
    is_liked       BOOLEAN                  DEFAULT FALSE,
    ignore         BOOLEAN                  DEFAULT FALSE,
    track_id       VARCHAR(100),
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建小红书笔记封面表
CREATE TABLE xhs_note_covers
(
    id         SERIAL PRIMARY KEY,
    note_id    INTEGER REFERENCES xhs_notes (id) ON DELETE CASCADE,
    url        TEXT,
    url_pre_path    TEXT,                           -- 将 'url_pre' 改为 'url_pre_path'
    url_default_path TEXT,                        -- 将 'url_default' 改为 'url_default_path'
    file_id    VARCHAR(100),
    width      INTEGER                  DEFAULT 0,
    height     INTEGER                  DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建小红书笔记图片表
CREATE TABLE xhs_note_images
(
    id         SERIAL PRIMARY KEY,
    note_id    INTEGER REFERENCES xhs_notes (id) ON DELETE CASCADE,
    url        TEXT,
    width      INTEGER                  DEFAULT 0,
    height     INTEGER                  DEFAULT 0,
    image_scene VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建小红书笔记交互信息表
CREATE TABLE xhs_note_interact_info
(
    id          SERIAL PRIMARY KEY,
    note_id     INTEGER REFERENCES xhs_notes (id) ON DELETE CASCADE,
    liked       BOOLEAN                  DEFAULT FALSE,
    liked_count INTEGER                  DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户收藏小红书笔记关系表
CREATE TABLE user_collect_xhs_notes
(
    id           SERIAL PRIMARY KEY,
    user_id      VARCHAR(50) REFERENCES users (uid),
    note_id      INTEGER REFERENCES xhs_notes (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, note_id)
);

-- 创建用户喜欢小红书笔记关系表
CREATE TABLE user_like_xhs_notes
(
    id           SERIAL PRIMARY KEY,
    user_id      VARCHAR(50) REFERENCES users (uid),
    note_id      INTEGER REFERENCES xhs_notes (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, note_id)
);

-- 创建用户收藏视频关系表
CREATE TABLE user_collect_videos
(
    id           SERIAL PRIMARY KEY,
    commenter_id VARCHAR(50) REFERENCES users (uid),
    video_id     INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (commenter_id, video_id)
);

-- 创建用户喜欢视频关系表
CREATE TABLE user_like_videos
(
    id           SERIAL PRIMARY KEY,
    commenter_id VARCHAR(50) REFERENCES users (uid),
    video_id     INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (commenter_id, video_id)
);

-- 创建子评论表
CREATE TABLE sub_comments
(
    id           SERIAL PRIMARY KEY,
    comment_id   VARCHAR(50) UNIQUE NOT NULL,
    parent_cmt_id VARCHAR(50) REFERENCES comments (comment_id) ON DELETE CASCADE,  -- 将 'parent_comment_id' 改为 'parent_cmt_id'
    commenter_id VARCHAR(50) REFERENCES users (uid),
    content      TEXT,
    ip_location  VARCHAR(100),
    user_digged  INTEGER                  DEFAULT 0,
    is_author_digged BOOLEAN             DEFAULT FALSE,
    is_hot       BOOLEAN                  DEFAULT FALSE,
    is_folded    BOOLEAN                  DEFAULT FALSE,
    user_buried  BOOLEAN                  DEFAULT FALSE,
    digg_count   INTEGER                  DEFAULT 0,
    create_time  BIGINT,
    last_modify_ts BIGINT,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sub_comments_parent_id ON sub_comments (parent_cmt_id);  -- 修改索引字段名
CREATE INDEX idx_sub_comments_commenter_id ON sub_comments (commenter_id);

-- 创建用户历史记录表
CREATE TABLE user_history_videos
(
    id           SERIAL PRIMARY KEY,
    commenter_id VARCHAR(50) REFERENCES users (uid),
    video_id     INTEGER REFERENCES videos (id) ON DELETE CASCADE,
    view_time    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (commenter_id, video_id)
);

-- 创建用户好友关系表
CREATE TABLE user_friends
(
    id        SERIAL PRIMARY KEY,
    user_id   INTEGER REFERENCES users (id) ON DELETE CASCADE,
    friend_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, friend_id)
);

-- 创建用户收藏帖子关系表
CREATE TABLE user_collect_posts
(
    id           SERIAL PRIMARY KEY,
    commenter_id VARCHAR(50) REFERENCES users (uid),
    post_id      INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (commenter_id, post_id)
);

-- 创建用户喜欢帖子关系表
CREATE TABLE user_like_posts
(
    id           SERIAL PRIMARY KEY,
    commenter_id VARCHAR(50) REFERENCES users (uid),
    post_id      INTEGER REFERENCES posts (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (commenter_id, post_id)
);

-- 创建索引
CREATE INDEX idx_videos_author_user_id ON videos (author_user_id);
CREATE INDEX idx_videos_type ON videos (video_type);             -- 这里索引名称与实际字段不匹配
CREATE INDEX idx_videos_aweme_id ON videos (aweme_id);
CREATE INDEX idx_comments_video_id ON comments (video_id);
CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_aweme_id ON comments (aweme_id);
CREATE INDEX idx_comments_user_id ON comments (commenter_id);     -- 这里索引名称与实际字段不匹配
CREATE INDEX idx_posts_author_user_id ON posts (author_user_id);
CREATE INDEX idx_posts_aweme_id ON posts (aweme_id);
CREATE INDEX idx_post_images_post_id ON post_images (post_id);
CREATE INDEX idx_post_videos_post_id ON post_videos (post_id);
CREATE INDEX idx_post_covers_post_id ON post_covers (post_id);
CREATE INDEX idx_post_status_post_id ON post_status (post_id);
CREATE INDEX idx_post_text_extra_post_id ON post_text_extra (post_id);
CREATE INDEX idx_post_share_info_post_id ON post_share_info (post_id);
CREATE INDEX idx_post_video_tags_post_id ON post_video_tags (post_id);
CREATE INDEX idx_user_collect_videos_user_id ON user_collect_videos (commenter_id);
CREATE INDEX idx_user_like_videos_user_id ON user_like_videos (commenter_id);
CREATE INDEX idx_user_history_videos_user_id ON user_history_videos (commenter_id);
CREATE INDEX idx_user_collect_music_user_id ON user_collect_music (commenter_id);
CREATE INDEX idx_user_collect_posts_user_id ON user_collect_posts (commenter_id);
CREATE INDEX idx_user_like_posts_user_id ON user_like_posts (commenter_id);

-- 小红书笔记相关索引
CREATE INDEX idx_xhs_notes_author_user_id ON xhs_notes (author_user_id);
CREATE INDEX idx_xhs_notes_note_id ON xhs_notes (note_id);
CREATE INDEX idx_xhs_notes_type ON xhs_notes (note_type);        -- 修改为与实际字段匹配
CREATE INDEX idx_xhs_note_covers_note_id ON xhs_note_covers (note_id);
CREATE INDEX idx_xhs_note_images_note_id ON xhs_note_images (note_id);
CREATE INDEX idx_xhs_note_interact_info_note_id ON xhs_note_interact_info (note_id);
CREATE INDEX idx_user_collect_xhs_notes_user_id ON user_collect_xhs_notes (user_id);
CREATE INDEX idx_user_like_xhs_notes_user_id ON user_like_xhs_notes (user_id);

CREATE INDEX idx_user_friends_user_id ON user_friends (user_id);
CREATE INDEX idx_user_friends_friend_id ON user_friends (friend_id);
