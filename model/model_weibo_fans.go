package model

import "time"

type WeiboFans struct {
	Uid            int64      `json:"uid" gorm:"TYPE:BIGINT(20);unique_index:uid_fans_idx"`
	FansId         int64      `json:"fans_id" gorm:"TYPE:BIGINT(20);unique_index:uid_fans_idx"`
	Like           bool       `json:"like"`
	LikeMe         bool       `json:"like_me"`
	FollowMe       bool       `json:"follow_me"`
	Following      bool       `json:"following"`
	AllowAllActMsg bool       `json:"allow_all_act_msg"`
	CreatedAt      *time.Time `json:"created_at" gorm:"TYPE:TIMESTAMP"`
}
