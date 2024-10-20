CREATE TABLE health.user
(
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid           VARCHAR(255) UNIQUE                                             NOT NULL,
    name          VARCHAR(100),
    system_avatar VARCHAR(255),
    custom_avatar VARCHAR(255),
    bind_id       VARCHAR(100),
    did           VARCHAR(255),
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP                             NOT NULL,
    last_login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP                             NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP                             NOT NULL,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    INDEX idx_did (did)
) CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;
