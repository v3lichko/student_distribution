package models

type Group struct {
	tableName struct{} `pg:"groups"`

	Number   int `json:"number" pg:"number"`
	Capacity int `json:"capacity" pg:"capacity"`
}
