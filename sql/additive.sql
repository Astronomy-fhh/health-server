CREATE TABLE health.additive
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    `desc`       VARCHAR(511) DEFAULT '',
    gb         VARCHAR(50)  DEFAULT '',
    category   VARCHAR(100)  DEFAULT '',
    tags       BLOB         DEFAULT NULL,
    image_url  VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
