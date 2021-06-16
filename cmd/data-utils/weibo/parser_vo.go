package weibo

import "github.com/juxuny/data-utils/model"

type FansVo struct {
	model.WeiboUser
	model.WeiboFans
}

type FriendshipsResp struct {
	Users                 []FansVo `json:"users"`
	HasFilteredAttentions bool     `json:"has_filtered_attentions"`
	NextCursor            int64    `json:"next_cursor"`
	PreviousCursor        int64    `json:"previous_cursor"`
	TotalNumber           int64    `json:"total_number"`
	UseSinkStrategy       bool     `json:"use_sink_stragety"`
	HasFilteredFans       bool     `json:"has_filtered_fans"`
	DisplayTotalNumber    int64    `json:"display_total_number"`
	Ok                    int64    `json:"ok"`
}
