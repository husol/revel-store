package dataservice

import (
	"husol.org/mypham/app/models"
	"github.com/jinzhu/gorm"
)

type InformationRepo struct {
	db *gorm.DB
}

type InformationResult struct {
	PostedUser string
	models.Information
}

func NewInformationRepo() *InformationRepo {
	return &InformationRepo{db:GetDBConnection()}
}

func (repo *InformationRepo) GetAll() []InformationResult {
	var results []InformationResult
	repo.db.Table("information").Order("updated_at DESC").Select("users.full_name AS posted_user, information.*").Joins("left join users on information.id_user = users.id").Find(&results)

	return results
}

func (repo *InformationRepo) GetActiveAll() []models.Information {
	var objects []models.Information
	repo.db.Where("status > 0").Order("updated_at DESC").Find(&objects)
	return objects
}

func (repo *InformationRepo) GetById(id int) *models.Information {
	var obj models.Information
	repo.db.First(&obj, id)
	return &obj
}

func (repo *InformationRepo) Create(obj *models.Information) error {
	tx := repo.db.Begin()

	if err := tx.Create(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *InformationRepo) Update(obj *models.Information) error {
	tx := repo.db.Begin()

	if err := tx.Save(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *InformationRepo) Delete(obj *models.Information) error {
	tx := repo.db.Begin()

	if err := tx.Delete(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}