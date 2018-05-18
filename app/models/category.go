package models

type Category struct {
	Id        			uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name				string	`gorm:"type:varchar(100);not null;unique"`
}
