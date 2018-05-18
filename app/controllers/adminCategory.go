package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/routes"
)

type AdminCategory struct {
	*revel.Controller
}

func (c AdminCategory) Index() revel.Result {
	//Count new comments
	repoComment := dataservice.NewCommentRepo()
	newCommentNum := repoComment.CountNewComments()
	c.ViewArgs["newCommentNum"] = newCommentNum

	repo := dataservice.NewCategoryRepo()
	var categories []models.Category
	categories, _ = repo.GetAll()
	c.ViewArgs["categories"] = categories

	return c.Render()
}

func (c AdminCategory) LoadForm() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repo := dataservice.NewCategoryRepo()
	var data interface{}

	data = models.Category{}
	if id > 0 {
		category := repo.GetById(id)
		data = category
	}

	husAjax := models.HusAjax{}
	content := husAjax.Fetch("AdminCategory/Form.html", data)
	husAjax.SetHTML("common_dialog", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c AdminCategory) Update() revel.Result {
	var id int
	c.Params.Bind(&id, "id_record")
	name := c.Params.Get("name")
	repo := dataservice.NewCategoryRepo()

	if id > 0 {
		category := repo.GetById(id)
		category.Name = name
		repo.Update(category);
		c.Flash.Success("Cập nhật Loại sản phẩm thành công.")
	} else {
		category := models.Category{Name: name}
		repo.Create(&category)
		c.Flash.Success("Thêm Loại sản phẩm thành công.")
	}

	return c.Redirect(routes.AdminCategory.Index())
}

func (c AdminCategory) Delete() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repo := dataservice.NewCategoryRepo()
	category := repo.GetById(id)

	if category.Id > 0 {
		repo.Delete(category)
		c.Flash.Success("Xóa Loại sản phẩm thành công.")
	} else {
		c.Flash.Error("Không thể thực hiện xóa Loại sản phẩm.")
	}
	return c.Redirect(routes.AdminCategory.Index())
}