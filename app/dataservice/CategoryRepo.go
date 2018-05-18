package dataservice

import (
	"husol.org/mypham/app/models"
	"github.com/jinzhu/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo() *CategoryRepo {
	return &CategoryRepo{db:GetDBConnection()}
}

func (repo *CategoryRepo) GetAll() ([]models.Category, error) {
	var objects []models.Category
	repo.db.Order("name", true).Find(&objects)
	return objects, nil
}

func (repo *CategoryRepo) GetById(id int) *models.Category {
	var obj models.Category
	repo.db.First(&obj, id)
	return &obj
}

func (repo *CategoryRepo) Create(obj *models.Category) error {
	tx := repo.db.Begin()

	if err := tx.Create(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *CategoryRepo) Update(obj *models.Category) error {
	tx := repo.db.Begin()

	if err := tx.Save(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *CategoryRepo) Delete(obj *models.Category) error {
	tx := repo.db.Begin()

	if err := tx.Delete(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}