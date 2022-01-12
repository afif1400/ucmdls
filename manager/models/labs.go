package models

import (
	"gorm.io/gorm"
)

type Lab struct {
	gorm.Model
	Name    string   `json:"name"`
	Branch  string   `json:"branch"`
	Systems []System `json:"systems"`
}
