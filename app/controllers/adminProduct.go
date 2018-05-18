package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/routes"
	"strconv"
	"path"
	"strings"
)

type AdminProduct struct {
	*revel.Controller
}

type Data struct {
	Categories interface{}
	*models.Product
}

func (c AdminProduct) Index() revel.Result {
	//Count new comments
	repoComment := dataservice.NewCommentRepo()
	newCommentNum := repoComment.CountNewComments()
	c.ViewArgs["newCommentNum"] = newCommentNum

	repo := dataservice.NewProductRepo()
	products := repo.GetAll()
	c.ViewArgs["products"] = products
	return c.Render()
}

func (c AdminProduct) LoadForm() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repo := dataservice.NewProductRepo()
	var data Data

	repoCate := dataservice.NewCategoryRepo()
	categories, _ := repoCate.GetAll()
	data.Categories = categories

	data.Product = &models.Product{}
	if id > 0 {
		data.Product = repo.GetById(id)
	}

	husAjax := models.HusAjax{}
	content := husAjax.Fetch("AdminProduct/Form.html", data)
	husAjax.SetHTML("common_dialog", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c AdminProduct) Update() revel.Result {
	var id int
	c.Params.Bind(&id, "id_record")
	var idCategory uint
	c.Params.Bind(&idCategory, "category")
	modelName := c.Params.Get("model_name")
	name := c.Params.Get("name")
	cover := c.Params.Files["cover"]
	shortDescription := c.Params.Get("short_description")
	description := c.Params.Get("description")
	cleanPrice := strings.Replace(c.Params.Get("price"), ",", "", -1)
	price, _ := strconv.ParseFloat(cleanPrice, 64)

	var isFeature int
	c.Params.Bind(&isFeature,"is_feature")
	var status int
	c.Params.Bind(&status,"status")
	repo := dataservice.NewProductRepo()

	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	if id > 0 {
		product := repo.GetById(id)
		product.IdCategory = idCategory
		product.ModelName = modelName
		product.Name = name
		product.ShortDescription = shortDescription
		product.Description = description
		product.Price = price
		product.IsFeature = isFeature
		product.Status = status

		if len(cover) > 0 {
			//Delete old cover
			hus.DeleteDirFile(product.Cover)

			//Upload and resize cover
			uploadDir := "/public/images/products"
			coverFile := models.HusFile{uploadDir, cover}
			coverName := "product_" + strconv.FormatUint(uint64(product.Id), 10) + path.Ext(coverFile.Files[0].Filename)
			uploadCover := coverFile.UploadFile(strings.ToLower(coverName));
			if uploadCover {
				coverFile.ThumbnailImage(coverName, 480, 360)
				product.Cover = uploadDir+"/"+coverName;
			}
		}
		repo.Update(product)

		c.Flash.Success("Cập nhật Sản phẩm thành công.")
	} else {
		product := models.Product{
			IdCategory: idCategory,
			ModelName: modelName,
			Name: name,
			ShortDescription: shortDescription,
			Description: description,
			Price: price,
			IsFeature: isFeature,
			Status: status,
		}
		repo.Create(&product)

		if len(cover) > 0 {
			//Upload and resize cover
			uploadDir := "/public/images/products"
			coverFile := models.HusFile{uploadDir, cover}
			coverName := "info_" + strconv.FormatUint(uint64(product.Id), 10) + path.Ext(coverFile.Files[0].Filename)
			uploadCover := coverFile.UploadFile(strings.ToLower(coverName));
			if uploadCover {
				coverFile.ThumbnailImage(coverName, 500, 315)
				product.Cover = uploadDir + "/" + coverName;
				repo.Update(&product)
			}
		}
		c.Flash.Success("Thêm Sản phẩm thành công.")
	}

	return c.Redirect(routes.AdminProduct.Index())
}

func (c AdminProduct) Delete() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	hus := models.Hus{}
	repo := dataservice.NewProductRepo()
	product := repo.GetById(id)

	if product.Id > 0 {
		//Delete cover
		hus.DeleteDirFile(product.Cover)
		repo.Delete(product)
		c.Flash.Success("Xóa Sản phẩm thành công.")
	} else {
		c.Flash.Error("Không thể thực hiện xóa Sản phẩm.")
	}
	return c.Redirect(routes.AdminProduct.Index())
}