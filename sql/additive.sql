create table health.additive
(
    id         bigint unsigned auto_increment
        primary key,
    name       varchar(100)                            not null,
    `desc`     varchar(1000) default ''                null,
    gb         varchar(50)   default ''                null,
    category   varchar(100)  default ''                null,
    tags       blob                                    null,
    image_url  varchar(255)  default ''                null,
    created_at timestamp     default CURRENT_TIMESTAMP null,
    updated_at timestamp     default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint name
        unique (name)
);

