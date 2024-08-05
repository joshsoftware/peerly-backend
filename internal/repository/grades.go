package repository

type GradesStorer interface {
}

type Grade struct {
	Id     int64  `db:"id"`
	Name   string `db:"name"`
	Points int64  `db:"points"`
}
