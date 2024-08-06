package repository

import "context"

type GradesStorer interface {
	ListGrades(ctx context.Context) (gradesList []Grade, err error)
}

type Grade struct {
	Id     int64  `db:"id"`
	Name   string `db:"name"`
	Points int64  `db:"points"`
}
