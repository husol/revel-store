package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Id        			uint 	`gorm:"primary_key;AUTO_INCREMENT"`
	Email				string	`gorm:"type:varchar(100);not null;unique"`
	FullName            string	`gorm:"not null"`
	Password            string
	Avatar				string
	Mobile				string	`gorm:"type:varchar(16)"`
	Address             string
	Role				int		`gorm:"type:tinyint(4);not null;default:1"`
	Token				string
	Status				int		`gorm:"type:tinyint(4);not null;default:0"`
	LastLogin			time.Time

	gorm.Model
}