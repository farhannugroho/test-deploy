package model

type UOM struct {
	Model
	Name string `json:"name" binding:"required"`
}
