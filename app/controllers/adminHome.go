package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/models"
	"strings"
	"husol.org/mypham/app/routes"
)

type AdminHome struct {
	*revel.Controller
}

func (c AdminHome) Index() revel.Result {
	//Count new comments
	repoComment := dataservice.NewCommentRepo()
	newCommentNum := repoComment.CountNewComments()
	c.ViewArgs["newCommentNum"] = newCommentNum

	repoTransaction := dataservice.NewTransactionRepo()
	//Get new transactions
	newTransactions := repoTransaction.GetByStatus(0)
	c.ViewArgs["newTransactions"] = newTransactions
	c.ViewArgs["newTransactionNum"] = len(newTransactions)
	//Count pending transactions
	pendingTransactionNum := len(repoTransaction.GetByStatus(1))
	c.ViewArgs["pendingTransactionNum"] = pendingTransactionNum
	//Count successful transactions
	successfulTransactionNum := len(repoTransaction.GetByStatus(2))
	c.ViewArgs["successfulTransactionNum"] = successfulTransactionNum

	return c.Render()
}

func (c AdminHome) Update() revel.Result {
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

	return c.Redirect(routes.AdminHome.Index())
}

func (c AdminHome) Sendmail() revel.Result {
	email := c.Params.Get("email")
	subject := c.Params.Get("subject")
	content := c.Params.Get("content")

	//Validate
	if !c.Validation.Required(subject).Ok {
		c.Flash.Error("Vui lòng nhập Tiêu đề")
		return c.Redirect(routes.AdminHome.Index())
	}
	if !c.Validation.Required(content).Ok {
		c.Flash.Error("Vui lòng nhập Nội dung")
		return c.Redirect(routes.AdminHome.Index())
	}

	husMail := models.HusMail{"noreply@myphamtayninh.com", "!@#123^%$"}
	//Supply information to send email
	if email == "" {
		repoUser := dataservice.NewUserRepo()
		users := repoUser.GetByRole(1)
		for index, user := range users {
			if index > 0 {
				email += ","
			}
			email += user.Email
		}
	}

	receivers := strings.Split(email, ",")

	//Send email
	err := husMail.SendMail(receivers, subject, content)

	if err != nil {
		c.Flash.Error("Gửi mail thông báo không thành công.")
	} else {
		c.Flash.Success("Gửi mail thông báo thành công.")
	}

	return c.Redirect(routes.AdminHome.Index())
}