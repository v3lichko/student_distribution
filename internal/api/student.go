package api

type Student struct {
	ISU         int    `json:"isu"`
	FullName    string `json:"full_name"`
	Telegram    string `json:"telegram"`
	Score       int    `json:"score"`
	GroupNumber *int   `json:"group_number"`
}
