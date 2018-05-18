package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
)

type Contact struct {
	*revel.Controller
}

func (c Contact) Index() revel.Result {
	return c.Render()
}

func (c Contact) SendMessage(contact_name, contact_email, contact_msg string) revel.Result {
	c.Validation.Required(contact_name).Message("Vui lòng nhập Họ và tên")
	c.Validation.Required(contact_email).Message("Vui lòng nhập Email")
	c.Validation.Email(contact_email).Message("Vui lòng nhập đúng chuẩn Email")
	c.Validation.Required(contact_msg).Message("Vui lòng nhập Tin nhắn")

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/contact")
	}

	//All data validated
	responseCaptcha := c.Params.Form.Get("g-recaptcha-response")
	secret, _ := revel.Config.String("CAPTCHA_SECRET_KEY")
	husCaptcha := models.HusCaptcha {Secret: secret}

	isValid := husCaptcha.Verify(responseCaptcha)

	if isValid {
		husMail := models.HusMail{"noreply@myphamtayninh.com", "!@#123^%$"}

		//Supply information to send email
		receivers:= []string{"store@myphamtayninh.com"}
		subject := "Liên hệ từ Mỹ Phẩm Tây Ninh"
		bodyMessage := "Lời nhắn từ khách hàng: "+contact_name+" <"+contact_email+">\n\n"+contact_msg
		//Send email
		err := husMail.SendMail(receivers, subject, bodyMessage)

		if err != nil {
			c.Flash.Error("Send message unsuccessfully.")
		} else {
			c.Flash.Success("Send message successfully.")
		}
	} else {
		c.Flash.Error("Vui lòng xác thực Captcha thành công để gửi tin.")
	}

	return c.Redirect("/contact")
}