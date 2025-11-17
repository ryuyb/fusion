INSERT INTO streaming_platforms (type, name, description, base_url, logo_url, enabled, priority, metadata, created_at,
                                 updated_at)
VALUES ('douyu', '斗鱼直播', '斗鱼直播平台', 'https://www.douyu.com', 'https://ir.douyu.com/images/LOGO.jpg', TRUE, 100, '{}'::jsonb, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       ('bilibili', '哔哩哔哩直播', '哔哩哔哩直播平台', 'https://live.bilibili.com', 'https://www.iconfont.cn/', TRUE, 100, '{}'::jsonb, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    ON CONFLICT (name) DO NOTHING;
