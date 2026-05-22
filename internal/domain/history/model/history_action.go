package model

import "time"

type Type string
type Table string

var (
	Add    Type = "add"
	Edit   Type = "edit"
	Sort   Type = "sort"
	Link   Type = "link"
	Unlink Type = "unlink"
)

var (
	Sk Table = "sk"
)

type HistoryAction struct {
	Type      Type
	Table     Table
	Id        int
	CreatedAt time.Time
}
