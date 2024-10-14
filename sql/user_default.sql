CREATE TABLE health.user_default
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,                     -- 自增ID
    name       VARCHAR(100) UNIQUE NOT NULL,                                   -- 唯一名称
    img        VARCHAR(100),                                                   -- 描述，使用 BLOB
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                            -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
);
