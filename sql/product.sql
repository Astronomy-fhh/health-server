CREATE TABLE product
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    barcode    VARCHAR(100) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    additives  BLOB,
    images     BLOB,
    other_desc VARCHAR(511),
    create_uid VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX      idx_barcode (barcode),         -- 创建条形码索引
    FULLTEXT   INDEX idx_name_fulltext (name) -- 创建名称全文索引
);
