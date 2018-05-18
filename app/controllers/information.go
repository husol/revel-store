package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/models"
)

type Information struct {
	*revel.Controller
}

func (c Information) Index() revel.Result {
	repo := dataservice.NewInformationRepo()
	var information []models.Information
	information = repo.GetActiveAll()
	c.ViewArgs["information"] = information

	return c.Render();
}
