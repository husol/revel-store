package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/dataservice"
	"math"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/routes"
)

type AdminTransaction struct {
	*revel.Controller
}

const TRANSACTION_PAGESIZE int = 10

func (c AdminTransaction) Index() revel.Result {
	//Count new comments
	repoComment := dataservice.NewCommentRepo()
	newCommentNum := repoComment.CountNewComments()
	c.ViewArgs["newCommentNum"] = newCommentNum

	//Calculate total page
	repoTransaction := dataservice.NewTransactionRepo()
	transactions := repoTransaction.Search("", -1, -1)
	totalPage := math.Ceil(float64(len(transactions)) / float64(TRANSACTION_PAGESIZE))
	c.ViewArgs["totalPage"] = int(totalPage)

	return c.Render()
}

func (c AdminTransaction) Paging() revel.Result {
	var page int
	c.Params.Bind(&page, "page")

	repoTransaction := dataservice.NewTransactionRepo()
	transactions := repoTransaction.Search("", page, TRANSACTION_PAGESIZE)

	husAjax := models.HusAjax{}
	content := husAjax.Fetch("AdminTransaction/list.html", transactions)
	husAjax.SetHTML("tableTransaction", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c AdminTransaction) Detail() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repoTransaction := dataservice.NewTransactionRepo()
	transaction := repoTransaction.GetById(id)

	repoOrder := dataservice.NewOrderRepo()
	orders := repoOrder.GetOrdersByTransaction(id)

	var data struct {
		Transaction *models.Transaction
		Orders	[]models.Order
	}

	data.Transaction = transaction
	data.Orders = orders

	husAjax := models.HusAjax{}
	content := husAjax.Fetch("AdminTransaction/detail.html", data)
	husAjax.SetHTML("divOrder", content)
	return c.RenderJSON(husAjax.OutData(true))
}

func (c AdminTransaction) Update() revel.Result {
	var id int
	c.Params.Bind(&id, "id")
	var status int
	c.Params.Bind(&status, "status")
	repo := dataservice.NewTransactionRepo()

	transaction := repo.GetById(id)

	if transaction.Id > 0 {
		if (status == -1) {
			//Delete orders belong to this transaction
			repoOrder := dataservice.NewOrderRepo()
			orders := repoOrder.GetOrdersByTransaction(int(transaction.Id))
			for _, order := range orders {
				repoOrder.Delete(&order)
			}
			//Delete transaction
			repo.Delete(transaction)
			c.Flash.Success("Hủy Giao dịch thành công.")
		} else {
			transaction.Status = int8(status)
			repo.Update(transaction)
			c.Flash.Success("Cập nhật Giao dịch thành công.")
		}
	} else {
		c.Flash.Error("Giao dịch không tồn tại.")
	}

	return c.Redirect(routes.AdminTransaction.Index())
}