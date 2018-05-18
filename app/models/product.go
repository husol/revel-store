package models

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	Id        			uint 		`gorm:"primary_key;AUTO_INCREMENT"`
	IdCategory			uint		`gorm:"not null"`
	ModelName  	 	   	string		`gorm:"type:varchar(64)"`
	Name	            string		`gorm:"type:varchar(128)"`
	Cover				string		`gorm:"type:varchar(255)"`
	ShortDescription	string		`gorm:"type:varchar(255)"`
	Description			string		`gorm:"type:text"`
	Price				float64		`gorm:"type:float;not null;default:0"`
	IsFeature			int			`gorm:"type:tinyint(4);not null;default:0"`
	CountView			int64		`gorm:"type:int(11);not null;default:0"`
	Status				int			`gorm:"not null;default:1"`

	gorm.Model
}