package dataservice

import (
	"husol.org/mypham/app/models"
	"github.com/jinzhu/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

type ProductResult struct {
	CategoryName string
	models.Product
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{db:GetDBConnection()}
}

func (repo *ProductRepo) GetAll() []ProductResult {
	var results []ProductResult
	repo.db.Table("products").Order("updated_at DESC").Select("categories.name AS category_name, products.*").Joins("join categories on products.id_category = categories.id").Find(&results)

	return results
}

func (repo *ProductRepo) Search(keyword string, pageIndex int, pageSize int) []models.Product {
	var objects []models.Product
	keyword = "%"+keyword+"%"
	offset := (pageIndex-1)*pageSize
	if pageIndex == -1 {
		offset = -1
	}

	repo.db.Limit(pageSize).Offset(offset).
		Table("products").
		Order("updated_at DESC").
		Where("status > 0 AND model_name LIKE ? OR name LIKE ? OR short_description LIKE ?", keyword, keyword, keyword).
		Find(&objects)

	return objects
}

func (repo *ProductRepo) GetFeaturedProducts() []models.Product {
	var objects []models.Product
	repo.db.Limit(4).Where("status > 0 AND is_feature = 1").Order("updated_at DESC").Find(&objects)
	return objects
}

func (repo *ProductRepo) GetRecentProducts() []models.Product {
	var objects []models.Product
	repo.db.Limit(4).Where("status > 0").Order("updated_at DESC").Find(&objects)
	return objects
}

func (repo *ProductRepo) GetMostViewedProducts() []models.Product {
	var objects []models.Product
	repo.db.Limit(8).Where("status > 0").Order("count_view DESC").Find(&objects)
	return objects
}

func (repo *ProductRepo) GetProductsByDestiny(destiny string, limit int) []models.Product {
	var objects []models.Product
	repo.db.Limit(limit).Where("status > 0 AND id_category = ? AND destiny LIKE ?", 1, "%"+destiny+"%").Order("is_feature DESC, updated_at DESC").Find(&objects)
	return objects
}

func (repo *ProductRepo) GetProductsByCategory(id int) []models.Product {
	var objects []models.Product
	repo.db.Limit(4).Where("status > 0 AND id_category = ?", id).Order("is_feature DESC, updated_at DESC").Find(&objects)
	return objects
}

func (repo *ProductRepo) GetById(id int) *models.Product {
	var obj models.Product
	repo.db.First(&obj, id)
	return &obj
}

func (repo *ProductRepo) Create(obj *models.Product) error {
	tx := repo.db.Begin()

	if err := tx.Create(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *ProductRepo) Update(obj *models.Product) error {
	tx := repo.db.Begin()

	if err := tx.Save(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *ProductRepo) Delete(obj *models.Product) error {
	tx := repo.db.Begin()

	if err := tx.Delete(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}