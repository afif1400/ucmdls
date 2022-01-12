package models

import "gorm.io/gorm"

type System struct {
	gorm.Model
	Name      string `json:"name"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	IpAddress string `json:"ip_address"`
	ROM       string `json:"rom"`
	RAM       string `json:"ram"`
	OS        string `json:"os"`
}
