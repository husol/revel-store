package models

type Order struct {
	Id        			uint 	`gorm:"primary_key;AUTO_INCREMENT"`
	IdTransaction		uint	`gorm:"not null"`
	IdProduct  	 	   	uint	`gorm:"not null"`
	Quantity			int64	`gorm:"type:int(11);not null;default:0"`
	Amount				float64	`gorm:"type: decimal(15,4); not null; default:0.0000"`

	Product				*Product
}