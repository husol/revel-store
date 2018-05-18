package controllers

import (
	"github.com/revel/revel"
)

type Help struct {
	*revel.Controller
}

func (c Help) Index() revel.Result {
	return c.Render()
}
