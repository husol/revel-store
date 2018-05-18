package dataservice

import (
	"github.com/jinzhu/gorm"
	"husol.org/mypham/app/models"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{db: GetDBConnection()}
}

func (repo *UserRepo) GetAll() ([]models.User, error) {
	var objects []models.User
	repo.db.Find(&objects)
	return objects, nil
}

func (repo *UserRepo) GetByEmail(email string) *models.User {
	var obj models.User
	repo.db.Where(&models.User{Email: email}).First(&obj)
	return &obj
}

func (repo *UserRepo) GetByEmailPassword(email, password string) *models.User {
	var obj models.User
	repo.db.Where(&models.User{Email: email, Password: password}).First(&obj)
	return &obj
}

func (repo *UserRepo) GetByRole(role int) []models.User {
	var objects []models.User
	repo.db.Where("role = ? AND status = ?", role, "1").Find(&objects)
	return objects
}

func (repo *UserRepo) GetById(id int) *models.User {
	var obj models.User
	repo.db.First(&obj, id)
	return &obj
}

func (repo *UserRepo) Create(obj *models.User) error {
	tx := repo.db.Begin()

	if err := tx.Create(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *UserRepo) Update(obj *models.User) error {
	tx := repo.db.Begin()

	if err := tx.Save(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *UserRepo) Delete(obj *models.User) error {
	tx := repo.db.Begin()

	if err := tx.Delete(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
