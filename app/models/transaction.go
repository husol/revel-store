package models

import "github.com/jinzhu/gorm"

type Transaction struct {
	Id        			uint 	`gorm:"primary_key;AUTO_INCREMENT"`
	IdUser				uint	`gorm:"not null"`
	Amount				float64	`gorm:"type: decimal(15,4); not null; default:0.0000"`
	ContactInfo			string	`gorm:"type: text; not null"`
	DeliverPlace		string 	`gorm:"type: varchar(255)"`
	Note				string	`gorm:"type: varchar(255)"`
	Status				int8	`gorm:"type: tinyint(4); not null; default: 0"`

	ContactName			string	`gorm:"-"`
	ContactEmail		string	`gorm:"-"`
	ContactMobile		string	`gorm:"-"`
	User				*User
	gorm.Model
}