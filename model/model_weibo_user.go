package model

import "time"

type WeiboUser struct {
	Id                int64      `json:"id" gorm:"BIGINT(21);PRIMARY_KEY;AUTO_INCREMENT"`
	IdStr             string     `json:"idstr"`
	Class             int64      `json:"class"`
	ScreenName        string     `json:"screen_name"`
	Name              string     `json:"name"`
	Province          string     `json:"province"`
	City              string     `json:"city"`
	Location          string     `json:"location"`
	Description       string     `json:"description"`
	Url               string     `json:"url"`
	ProfileImageUrl   string     `json:"profile_image_url"`
	ProfileUrl        string     `json:"profile_url"`
	Domain            string     `json:"domain"`
	WeiHao            string     `json:"weihao"`
	Gender            string     `json:"gender"`
	FollowersCount    int64      `json:"followers_count"`
	FriendsCount      int64      `json:"friends_count"`
	PageFriendsCount  int64      `json:"pagefriends_count"`
	StatusesCount     int64      `json:"statuses_count"`
	VideoStatusCount  int64      `json:"video_status_count"`
	VideoPlayCount    int64      `json:"video_play_count"`
	FavouritesCount   int64      `json:"favourites_count"`
	CreatedAt         *time.Time `json:"created_at"`
	GeoEnabled        bool       `json:"geo_enabled"`
	Verified          bool       `json:"verified"`
	VerifiedType      int        `json:"verified_type"`
	Remark            string     `json:"remark"`
	StatusId          int64      `json:"status_id"`
	StatusIdstar      int64      `json:"status_idstar"`
	PType             int64      `json:"ptype"`
	AllowAllComment   bool       `json:"allow_all_comment"`
	AvatarLarge       string     `json:"avatar_large"`
	AvatarHd          string     `json:"avatar_hd"`
	VerifiedReason    string     `json:"verified_reason"`
	VerifiedTrade     string     `json:"verified_trade"`
	VerifiedReasonUrl string     `json:"verified_reason_url"`
	VerifiedSourceUrl string     `json:"verified_source_url"`
	OnlineStatus      int64      `json:"online_status"`
	BiFollowersCount  int64      `json:"bi_followers_count"`
	Lang              string     `json:"lang"`
	Star              int64      `json:"star"`
	Mbtype            int64      `json:"mbtype"`
	Mbrank            int64      `json:"mbrank"`
	Svip              int64      `json:"svip"`
	BlockWorld        int64      `json:"block_world"`
	BlockApp          int64      `json:"block_app"`
	CreditScore       int64      `json:"credit_score"`
	UserAbility       int64      `json:"user_ability"`
	Urank             int64      `json:"urank"`
	StoryReadState    int64      `json:"story_read_state"`
	VclubMember       int64      `json:"vclub_member"`
	IsTeenager        int64      `json:"is_teenager"`
	IsGuardian        int64      `json:"is_guardian"`
	IsTeenagerList    int64      `json:"is_teenager_list"`
	PcNew             int64      `json:"pc_new"`
	SpecialFollow     bool       `json:"special_follow"`
	PlanetVideo       int64      `json:"planet_video"`
	VideoMark         int64      `json:"video_mark"`
	LiveStatus        int64      `json:"live_status"`
	UserAbilityExtend int64      `json:"user_ability_extend"`
	BrandAccount      int64      `json:"brand_account"`
}
