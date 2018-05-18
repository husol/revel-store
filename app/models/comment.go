package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	Id        			uint 	`gorm:"primary_key;AUTO_INCREMENT"`
	IdUser        		uint 	`gorm:"not null"`
	IdProduct        	uint 	`gorm:"not null"`
	Content            	string	`gorm:"not null"`
	Status				int		`gorm:"type:tinyint(4);not null;default:0"`

	User				*User
	Product				*Product
	gorm.Model
}