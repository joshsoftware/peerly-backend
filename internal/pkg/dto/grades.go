package dto

type Grade struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Points int64  `json:"points"`
}

type UpdateGradeReq struct {
	Points int64 `json:"points"`
	Id     int64
}
