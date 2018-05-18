package models

import (
	"github.com/jinzhu/gorm"
)

type Information struct {
	Id        			uint 	`gorm:"primary_key;AUTO_INCREMENT"`
	IdUser				uint	`gorm:"not null"`
	Title				string	`gorm:"type:varchar(255);not null"`
	Cover				string
	Description			string	`gorm:"type:text"`
	Status				int		`gorm:"type:tinyint(4);not null;default:0"`
	gorm.Model
}