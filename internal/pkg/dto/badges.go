package dto

type Badge struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	RewardPoints int64  `json:"reward_points"`
	UpdatedBy    string `json:"updated_by"`
}

type UpdateBadgeReq struct {
	RewardPoints int64 `json:"reward_points"`
	Id           int64
	UserId       int64
}
