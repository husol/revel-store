package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/dataservice"
	"os"
	"log"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	repoCate := dataservice.NewCategoryRepo()
	categories, _ := repoCate.GetAll()
	c.ViewArgs["categories"] = categories

	repoProduct := dataservice.NewProductRepo()

	featuredProducts := repoProduct.GetFeaturedProducts()
	c.ViewArgs["featuredProducts"] = featuredProducts

	mostViewedProducts := repoProduct.GetMostViewedProducts()
	c.ViewArgs["mostViewedProducts"] = mostViewedProducts

	return c.Render()
}

func (c App) Sitemap() revel.Result {
	file, err := os.Open(revel.BasePath + "/public/sitemap.xml")
	if err != nil {
		log.Fatal("Cannot open sitemap.xml file", err)
	}

	return c.RenderFile(file, revel.Inline)
}

func (c App) Robots() revel.Result {
	file, err := os.Open(revel.BasePath + "/public/robots.txt")
	if err != nil {
		log.Fatal("Cannot open robots.txt file", err)
	}

	return c.RenderFile(file, revel.Inline)
}
