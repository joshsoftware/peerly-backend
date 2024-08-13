package dto

type Badge struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	RewardPoints int64  `json:"reward_points"`
}

type UpdateBadgeReq struct {
	RewardPoints int64 `json:"reward_points"`
	Id           int64
}
