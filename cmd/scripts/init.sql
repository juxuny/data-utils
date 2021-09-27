create schema split_video collate utf8mb4_general_ci charset 'utf8mb4';
CREATE TABLE eng_movie (
    id int(11) auto_increment primary key,
    name varchar(100) not null,
    parent_id int(11) not null default 0,
    create_time timestamp not null default current_timestamp,
    index(name),
    index(parent_id),
    index(create_time)
) comment '电影';

CREATE TABLE eng_subtitle (
    id int(11) auto_increment primary key ,
    movie_id int(11) not null,
    ext varchar(20) not null,
    file_name varchar(200) not null,
    create_time timestamp not null default current_timestamp,
    index(movie_id),
    index(file_name),
    index(create_time)
) comment '字幕';

CREATE TABLE eng_subtitle_block (
    id int(11) not null primary key auto_increment,
    subtitle_id int(11) not null,
    block_id int(11) not null,
    start_time varchar(20) not null,
    end_time varchar(20) not null,
    duration_extend varchar(200),
    content text not null,
    create_time timestamp not null default  current_timestamp,
    index(block_id),
    index(start_time),
    index(end_time),
    index(create_time),
    index(subtitle_id)
) comment 'subtitle block data';