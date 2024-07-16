package dto

type Reward struct {
	Id             int64 `json:"id"`
	AppreciationId int64 `json:"appreciation_id"`
	Point          int64 `json:"point"`
	SenderId       int64 `json:"sender"`
	CreatedAt      int64 `json:"created_at"`
}
