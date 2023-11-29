package dao

import "gorm.io/gorm"

type RuleGroupData struct {
	gorm.Model
	Name        string
	RuleNumber  int
	Description string
}
