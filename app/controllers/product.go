package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/models"
	"math"
)

type Product struct {
	*revel.Controller
}

func (c Product) Index() revel.Result {
	//Get params if any
	keyword := c.Params.Get("keyword")
	var page int
	c.Params.Bind(&page, "page")

	repoCate := dataservice.NewCategoryRepo()
	categories, _ := repoCate.GetAll()
	c.ViewArgs["categories"] = categories

	repoProduct := dataservice.NewProductRepo()

	recentProducts := repoProduct.GetRecentProducts()
	c.ViewArgs["recentProducts"] = recentProducts

	otherProductsTotal := len(repoProduct.Search(keyword, -1, -1))
	pagesTotal := int(math.Ceil(float64(otherProductsTotal)/12))
	c.ViewArgs["pagesTotal"] = pagesTotal

	if page < 1 {
		page = 1
	} else if page > pagesTotal {
		page = pagesTotal
	}

	otherProducts := repoProduct.Search(keyword, page, 12)
	c.ViewArgs["otherProducts"] = otherProducts

	c.ViewArgs["currpage"] = page
	c.ViewArgs["keyword"] = keyword

	return c.Render()
}

func (c Product) Category() revel.Result {
	var idCate int
	c.Params.Bind(&idCate, "id")

	repoCate := dataservice.NewCategoryRepo()

	//Get my category
	myCategory := repoCate.GetById(idCate)
	c.ViewArgs["myCategory"] = myCategory

	//Get categories list
	categories, _ := repoCate.GetAll()
	c.ViewArgs["categories"] = categories

	repoProduct := dataservice.NewProductRepo()

	//Get featured products for slider
	featuredProducts := repoProduct.GetFeaturedProducts()
	c.ViewArgs["featuredProducts"] = featuredProducts

	//Get products by category
	products := repoProduct.GetProductsByCategory(idCate)
	c.ViewArgs["products"] = products

	return c.Render()
}

func (c Product) Detail() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repoProduct := dataservice.NewProductRepo()
	myProduct := repoProduct.GetById(id)

	if myProduct.Id == 0 {
		return c.NotFound("The product is not existed.")
	}

	//Get logged user
	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)
	c.ViewArgs["LoggedUser"] = loggedUser

	//Get comments of the product
	repoComment := dataservice.NewCommentRepo()
	comments := repoComment.GetCommentsByProduct(int(myProduct.Id))
	c.ViewArgs["Comments"] = comments

	repoCate := dataservice.NewCategoryRepo()
	myCategory := repoCate.GetById(int(myProduct.IdCategory))
	c.ViewArgs["myCate"] = myCategory

	//Increase countview of product
	myProduct.CountView += 1
	repoProduct.Update(myProduct)

	c.ViewArgs["myProduct"] = myProduct

	//Get related products by category
	relatedProducts := repoProduct.GetProductsByCategory(int(myCategory.Id))
	c.ViewArgs["relatedProducts"] = relatedProducts

	return c.Render()
}

func (c Product) UpdateComment() revel.Result {
	var id int
	c.Params.Bind(&id, "id")
	var id_product int
	c.Params.Bind(&id_product, "id_product")
	content := c.Params.Get("content")

	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	if loggedUser.Id == 0 {
		return c.RenderJSON("expired_session")
	}

	repoComment := dataservice.NewCommentRepo()

	husAjax := models.HusAjax{}
	if id > 0 {
		myComment := repoComment.GetById(id)
		if myComment.Id == 0 {
			return c.RenderJSON(husAjax.OutData(false))
		}
		myComment.Content = content;
		myComment.Status = 0;
		repoComment.Update(myComment)
	} else {
		myComment := models.Comment{IdUser: loggedUser.Id, IdProduct:uint(id_product), Content:content}
		repoComment.Create(&myComment)
	}

	var data struct {
		LoggedUser models.User
		Comments interface{}
	}

	data.LoggedUser = loggedUser
	data.Comments = repoComment.GetCommentsByProduct(id_product)

	htmlContent := husAjax.Fetch("Product/comments.html", data)
	husAjax.SetHTML("comments", htmlContent)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c Product) DeleteComment() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	if loggedUser.Id == 0 {
		return c.RenderJSON("expired_session")
	}

	repoComment := dataservice.NewCommentRepo()
	myComment := repoComment.GetById(id)

	husAjax := models.HusAjax{}
	if myComment.Id == 0 {
		return c.RenderJSON(husAjax.OutData(false))
	}

	repoComment.Delete(myComment)
	return c.RenderJSON(husAjax.OutData(myComment))
}