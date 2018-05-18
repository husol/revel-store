package dataservice

import (
	"husol.org/mypham/app/models"
	"github.com/jinzhu/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{db:GetDBConnection()}
}

func (repo *OrderRepo) GetOrdersByTransaction(id int) []models.Order {
	var objects []models.Order
	repo.db.Where("id_transaction = ?", id).Find(&objects)
	productRepo := NewProductRepo()
	for index := range objects {
		objects[index].Product = productRepo.GetById(int(objects[index].IdProduct))
	}
	return objects
}

func (repo *OrderRepo) GetById(id int) *models.Order {
	var obj models.Order
	repo.db.First(&obj, id)
	return &obj
}

func (repo *OrderRepo) Create(obj *models.Order) error {
	tx := repo.db.Begin()

	if err := tx.Create(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *OrderRepo) Update(obj *models.Order) error {
	tx := repo.db.Begin()

	if err := tx.Save(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *OrderRepo) Delete(obj *models.Order) error {
	tx := repo.db.Begin()

	if err := tx.Delete(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}