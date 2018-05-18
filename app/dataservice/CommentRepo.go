package dataservice

import (
	"husol.org/mypham/app/models"
	"github.com/jinzhu/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

type CommentResult struct {
	UserFullName 		string
	ProductModelName 	string
	ProductName 		string
	models.Comment
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{db:GetDBConnection()}
}

func (repo *CommentRepo) GetActiveAll() []models.Comment {
	var objects []models.Comment
	repo.db.Where("status > 0").Order("updated_at DESC").Find(&objects)
	return objects
}

func (repo *CommentRepo) CountNewComments() int {
	var count int
	repo.db.Table("comments").Where("status = 0").Count(&count)
	return count
}

func (repo *CommentRepo) GetCommentsByProduct(id int) []CommentResult {
	var results []CommentResult
	repo.db.Table("comments").
	Order("created_at DESC").
	Where("comments.id_product = ?", id).
	Joins("join users on comments.id_user = users.id").
	Joins("join products on comments.id_product = products.id").
	Select("users.full_name AS user_full_name, products.model_name AS product_model_name, products.name AS product_name, comments.*").
	Find(&results)
	return results
}

func (repo *CommentRepo) Search(keyword string, pageIndex int, pageSize int) []models.Comment {
	var objects []models.Comment
	keyword = "%"+keyword+"%"
	offset := (pageIndex-1)*pageSize
	if pageIndex == -1 {
		offset = -1
	}

	repo.db.Limit(pageSize).Offset(offset).
		Table("comments").
		Order("status ASC, updated_at DESC").
		Joins("JOIN users on comments.id_user = users.id").
		Joins("JOIN products on comments.id_product = products.id").
		Where("users.full_name LIKE ? OR products.model_name LIKE ? OR products.name LIKE ?", keyword, keyword, keyword).
		Find(&objects)

	repoUser := NewUserRepo()
	repoProduct := NewProductRepo()
	for index := range objects {
		objects[index].User = repoUser.GetById(int(objects[index].IdUser))
		objects[index].Product = repoProduct.GetById(int(objects[index].IdProduct))
	}
	return objects
}

func (repo *CommentRepo) GetById(id int) *models.Comment {
	var obj models.Comment
	repo.db.First(&obj, id)
	return &obj
}

func (repo *CommentRepo) Create(obj *models.Comment) error {
	tx := repo.db.Begin()

	if err := tx.Create(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *CommentRepo) Update(obj *models.Comment) error {
	tx := repo.db.Begin()

	if err := tx.Save(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *CommentRepo) Delete(obj *models.Comment) error {
	tx := repo.db.Begin()

	if err := tx.Delete(obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}