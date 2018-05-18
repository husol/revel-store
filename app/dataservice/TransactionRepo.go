package dataservice

import (
	"husol.org/mypham/app/models"
	"github.com/jinzhu/gorm"
	"encoding/json"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo() *TransactionRepo {
	return &TransactionRepo{db:GetDBConnection()}
}

func (repo *TransactionRepo) Search(keyword string, pageIndex int, pageSize int) []models.Transaction {
	var objects []models.Transaction
	keyword = "%"+keyword+"%"
	offset := (pageIndex-1)*pageSize
	if pageIndex == -1 {
		offset = -1
	}

	repo.db.Limit(pageSize).Offset(offset).
		Table("transactions").
		Order("status ASC, updated_at ASC").
		Joins("LEFT JOIN users on transactions.id_user = users.id").
		Where("transactions.id LIKE ? OR users.full_name LIKE ?", keyword, keyword).
		Find(&objects)

	repoUser := NewUserRepo()
	var contact struct {
		Name	string
		Email	string
		Mobile	string
	}
	for index := range objects {
		json.Unmarshal([]byte(objects[index].ContactInfo), &contact)
		objects[index].ContactName = contact.Name
		objects[index].ContactEmail = contact.Email
		objects[index].ContactMobile = contact.Mobile
		objects[index].User = repoUser.GetById(int(objects[index].IdUser))
	}
	return objects
}

func (repo *TransactionRepo) GetByStatus(status int) []models.Transaction {
	var objects []models.Transaction
	repo.db.Where("status = ?", status).Find(&objects)

	repoUser := NewUserRepo()
	var contact struct {
		Name	string
		Email	string
		Mobile	string
	}
	for index := range objects {
		json.Unmarshal([]byte(objects[index].ContactInfo), &contact)
		objects[index].ContactName = contact.Name
		objects[index].ContactEmail = contact.Email
		objects[index].ContactMobile = contact.Mobile
		objects[index].User = repoUser.GetById(int(objects[index].IdUser))
	}
	return objects
}

func (repo *TransactionRepo) GetById(id int) *models.Transaction {
	var obj models.Transaction
	repo.db.First(&obj, id)
	return &obj
}

func (repo *TransactionRepo) Create(obj *models.Transaction) error {
	tx := repo.db.Begin()

	if err := tx.Create(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *TransactionRepo) Update(obj *models.Transaction) error {
	tx := repo.db.Begin()

	if err := tx.Save(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *TransactionRepo) Delete(obj *models.Transaction) error {
	tx := repo.db.Begin()

	if err := tx.Delete(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}