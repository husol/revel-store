package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
	"encoding/json"
	"husol.org/mypham/app/dataservice"
)

type Transaction struct {
	*revel.Controller
}

func (c Transaction) Index() revel.Result {

	return c.Render();
}

func (c Transaction) Checkout() revel.Result {
	responseCaptcha := c.Params.Get("g-recaptcha-response")
	secret, _ := revel.Config.String("CAPTCHA_SECRET_KEY")
	husCaptcha := models.HusCaptcha {Secret: secret}

	husAjax := models.HusAjax{}

	isValid := husCaptcha.Verify(responseCaptcha)
	if !isValid {
		husAjax.SetMessage("Vui lòng xác thực Captcha thành công.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	//Parse contactInfo data
	contactInfo := c.Params.Get("contactInfo")
	var contact struct {
		Name	string
		Email	string
		Mobile	string
	}
	json.Unmarshal([]byte(contactInfo), &contact)

	//Parse trolley data
	deliverPlace := c.Params.Get("deliverPlace")
	note := c.Params.Get("note")
	trolleyString := c.Params.Get("trolley")
	var trolley map[int] struct {
		Model		string
		Quantity	int64
		Price		float64
		Amount		float64
	}
	json.Unmarshal([]byte(trolleyString), &trolley)
	//Get loggedUser if any
	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	//Create transaction
	var transaction models.Transaction
	transaction = models.Transaction{IdUser: loggedUser.Id, ContactInfo: contactInfo, DeliverPlace: deliverPlace, Note: note}
	repoTransaction := dataservice.NewTransactionRepo()
	repoTransaction.Create(&transaction)

	repoOrder := dataservice.NewOrderRepo()
	total := 0.0;
	for id, row := range trolley {
		//Create order
		order := models.Order{IdTransaction: transaction.Id, IdProduct: uint(id), Quantity: row.Quantity, Amount: row.Amount}
		repoOrder.Create(&order)
		total += order.Amount
	}

	transaction.Amount = total
	repoTransaction.Update(&transaction)

	transaction.ContactName = contact.Name
	transaction.ContactEmail = contact.Email
	transaction.ContactMobile = contact.Mobile

	//Send notification mails
	husMail := models.HusMail{"noreply@myphamtayninh.com", "!@#123^%$"}

	var data struct {
		Trans *models.Transaction
		Orders []models.Order
	}
	data.Trans = &transaction
	data.Orders = repoOrder.GetOrdersByTransaction(int(transaction.Id))
	bodyMessage := husAjax.Fetch("Transaction/invoice.html", data)

	//Supply information to send client email
	receivers:= []string{contact.Email}
	subject := "Đặt hàng thành công tại Mỹ Phẩm Tây Ninh"

	err := husMail.SendMail(receivers, subject, bodyMessage)
	if err != nil {
		husAjax.SetMessage("Gửi mail đơn hàng không thành công.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	//Supply information to send admin email
	receivers= []string{"store@myphamtayninh.com"}
	subject = "Đơn hàng mới tại Mỹ Phẩm Tây Ninh"

	err = husMail.SendMail(receivers, subject, bodyMessage)
	if err != nil {
		husAjax.SetMessage("Gửi mail đơn hàng không thành công.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	return c.RenderJSON(husAjax.OutData(true))
}
