package model

type SkUri struct {
	Id int `uri:"id" binding:"required,min=1,sk"`
}

type SkAdd struct {
	Id       int    `json:"id" binding:"required,min=1"`
	Name     string `json:"name" binding:"required"`
	IsActive bool   `json:"isActive"`
}

type SkEdit struct {
	Name     string `json:"name" binding:"required"`
	IsActive bool   `json:"isActive"`
}

type SkItem struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type SkFilter struct {
	Offset int `form:"offset" binding:"min=0"`
	Limit  int `form:"limit" binding:"required,min=1"`
}
