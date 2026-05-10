package models

type GroupDistribution struct {
	GroupNumber int       `json:"group_number"`
	Students    []Student `json:"students"`
}
