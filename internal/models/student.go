package models

type Student struct {
	tableName struct{} `pg:"students"`

	ISU         int    `json:"isu" pg:"isu"`
	FullName    string `json:"full_name" pg:"full_name"`
	Telegram    string `json:"telegram" pg:"telegram"`
	Score       int    `json:"score" pg:"score"`
	GroupNumber *int   `json:"group_number" pg:"group_number"`
}
