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