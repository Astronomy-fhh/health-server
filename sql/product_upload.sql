CREATE TABLE health.product_upload
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,                     -- 自增ID
    barcode    VARCHAR(100) NOT NULL,                                          -- 条形码
    name       VARCHAR(100) NOT NULL,                                          -- 名称
    additives  BLOB,                                                           -- 添加剂，存储为 BLOB
    images     BLOB,                                                           -- 图片，存储为 BLOB
    other_desc VARCHAR(511),                                                   -- 描述，最大511字符
    create_uid VARCHAR(100) NOT NULL,                                          -- 创建用户ID
    stats      int(10) DEFAULT 0,                                               -- 状态，0: 待审核，1: 审核通过，2: 审核不通过
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                            -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
);
