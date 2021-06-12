CREATE SCHEMA crawl charset 'utf8mb4' collate utf8mb4_general_ci;

DROP TABLE IF EXISTS job;
CREATE TABLE job (
    id bigint(21) primary key auto_increment,
    job_type varchar(60) not null,
    state varchar(20) not null comment '任务状态',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp null on update current_timestamp,
    started_at timestamp null comment '任务执行开始时间',
    end_at timestamp null comment '任务执行结束时间',
    meta_data text comment '任务内容，一般是JSON',
    result text comment '执行结果',
    INDEX(job_type),
    INDEX(state),
    INDEX(created_at),
    INDEX(started_at),
    INDEX(end_at)
) comment 'job queue'


DROP TABLE IF EXISTS ad_email;
CREATE TABLE ad_email (
    id bigint(21) primary key auto_increment,
    email varchar(64) not null unique ,
    count int(11) not null default 0,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp null,
    last_error TEXT,
    INDEX(created_at)
);

DROP TABLE IF EXISTS proxy;
CREATE TABLE proxy
(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ip VARCHAR(20) NOT NULL,
    port INT NOT NULL,
    provider_name VARCHAR(60) NULL,
    socks5 TINYINT NULL,
    http TINYINT NULL,
    https TINYINT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    latency DECIMAL(20,2) NULL,
    used_count INT NULL,
    UNIQUE INDEX(ip, port),
    INDEX(created_at),
    INDEX(updated_at),
    INDEX(latency),
    INDEX(used_count)
);
