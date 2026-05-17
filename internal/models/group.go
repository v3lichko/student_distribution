package models

type Group struct {
	tableName struct{} `pg:"groups"`

	Number   int `pg:"number"`
	Capacity int `pg:"capacity"`
}
