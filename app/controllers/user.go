package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/dataservice"
	"encoding/json"
	"crypto/md5"
	"io"
	"encoding/hex"
	"strconv"
	"path"
	"strings"
	"time"
	"husol.org/mypham/app/routes"
)

type User struct {
	*revel.Controller
}

func (c User) Index() revel.Result {
	email := c.Params.Get("email")
	token := c.Params.Get("token")

	//Check token
	repo := dataservice.NewUserRepo()

	user := repo.GetByEmail(email)

	if user.Id > 0 {
		if user.Token != token {
			c.Flash.Error("Kích hoạt không thành công.")
		} else {
			user.Token = ""
			user.Status = 1
			repo.Update(user)
			c.Flash.Success("Bạn đã kích hoạt tài khoản thành công.")
		}
	} else {
		c.Flash.Error("Tài khoản không tồn tại trong hệ thống.")
	}

	return  c.Redirect(routes.App.Index())
}

func (c User) Login() revel.Result {
	user := models.User{}
	json.Unmarshal([]byte(c.Session["loggedUser"]), &user)

	if user.Id > 0 {
		return c.Redirect(routes.App.Index())
	}

	return  c.Render()
}

func (c User) PostLogin() revel.Result {
	email := c.Params.Get("email")
	password := c.Params.Get("password")

	c.Validation.Required(email).Message("Vui lòng nhập Email")
	c.Validation.Email(email).Message("Vui lòng nhập đúng chuẩn Email")
	c.Validation.Required(password).Message("Vui lòng nhập Password")
	c.Validation.MinSize(password, 8).Message("Password có ít nhất 8 ký tự")

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Login())
	}
	//All data validated

	//Check Login
	repo := dataservice.NewUserRepo()

	//Encrypt md5 password
	h := md5.New()
	io.WriteString(h, password)
	encryptPassword := hex.EncodeToString(h.Sum(nil))

	user := repo.GetByEmailPassword(email, encryptPassword)

	if user.Id > 0 {
		if user.Status == 0 {
			c.Flash.Error("Vui lòng kích hoạt tài khoản của bạn.")
		} else {
			user.LastLogin = time.Now()
			repo.Update(user)
			jsonUser, _ := json.Marshal(user)
			c.Session["loggedUser"] = string(jsonUser)

			c.Flash.Success("Chào "+ user.FullName + " đến với Mỹ Phẩm Tây Ninh")
			if user.Role == 0 {
				return c.Redirect(routes.AdminHome.Index())
			}

			return c.Redirect(routes.App.Index())
		}

	} else {
		c.Flash.Error("Invalid Email or Password")
	}

	return  c.Redirect(routes.User.Login())
}

func (c User) Logout() revel.Result {
	delete(c.Session, "loggedUser")
	c.Flash.Success("Đăng xuất thành công.")

	return  c.Redirect(routes.App.Index())
}

func (c User) RegisterUpdate() revel.Result {
	fullname := c.Params.Get("fullname")
	email := c.Params.Get("email")
	password := c.Params.Get("password")
	confirmPassword := c.Params.Get("confirmpassword")
	mobile := c.Params.Get("mobile")
	address := c.Params.Get("address")

	husAjax := models.HusAjax{}

	//Validate
	if !c.Validation.Required(fullname).Ok {
		husAjax.SetMessage("Vui lòng nhập Họ và tên")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Required(email).Ok {
		husAjax.SetMessage("Vui lòng nhập Email")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Email(email).Ok {
		husAjax.SetMessage("Vui lòng nhập đúng chuẩn Email")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Required(password).Ok {
		husAjax.SetMessage("Vui lòng nhập Password")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.MinSize(password, 8).Ok {
		husAjax.SetMessage("Password có ít nhất 8 ký tự")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if password != confirmPassword {
		husAjax.SetMessage("Nhập lại mật khẩu không trùng khớp mật khẩu đã nhập")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Required(mobile).Ok {
		husAjax.SetMessage("Vui lòng nhập Số điện thoại")
		return c.RenderJSON(husAjax.OutData(false))
	}
	//All data validated

	responseCaptcha := c.Params.Form.Get("g-recaptcha-response")
	secret, _ := revel.Config.String("CAPTCHA_SECRET_KEY")

	husCaptcha := models.HusCaptcha {Secret: secret}

	isValid := husCaptcha.Verify(responseCaptcha)

	if !isValid {
		husAjax.SetMessage("Vui lòng xác thực Captcha thành công.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	repo := dataservice.NewUserRepo()

	//Check user if exist
	user := repo.GetByEmail(email)
	if user.Id > 0 && user.Status == 1 {
		husAjax.SetMessage("Email đã tồn tại trong hệ thống.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	//Encrypt md5 password
	h := md5.New()
	io.WriteString(h, password)
	encryptPassword := hex.EncodeToString(h.Sum(nil))

	//Generate Token
	hus := models.Hus{}
	token := hus.RandomString(16)

	husMail := models.HusMail{"noreply@myphamtayninh.com", "!@#123^%$"}

	//Supply information to send email
	receivers:= []string{email}
	subject := "Kích hoạt tài khoản Mỹ Phẩm Tây Ninh"
	bodyMessage := "Chào mừng quý khách đến với Mỹ Phẩm Tây Ninh.<br/>Vui lòng click vào <a href=\""+hus.BaseUrl()+"/user/"+email+"/"+token+"\">đây</a> để kích hoạt tài khoản của bạn."
	//Send email
	err := husMail.SendMail(receivers, subject, bodyMessage)

	if err != nil {
		husAjax.SetMessage("Gửi mail xác thực không thành công.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	if user.Id > 0 {
		user.FullName = fullname
		user.Password = encryptPassword
		user.Mobile = mobile
		user.Address = address
		user.Token = token
		repo.Update(user)
	} else {
		user := models.User{FullName: fullname, Email:email, Password:encryptPassword, Mobile:mobile, Address:address, Token:token}
		repo.Create(&user)
	}

	c.Flash.Success("Vui lòng kích hoạt tài khoản từ địa chỉ email đã đăng ký.")

	return c.RenderJSON(husAjax.OutData(true))
}

func (c User) AccountForm(typeForm int) revel.Result {
	hus := models.Hus{}
	husAjax := models.HusAjax{}

	loggedUser := models.User{}
	if typeForm == 1 {
		hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)
	}

	content := husAjax.Fetch("User/AccountForm.html", loggedUser)
	husAjax.SetHTML("common_dialog", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c User) AccountUpdate() revel.Result {
	var id int
	c.Params.Bind(&id, "id_record")
	fullname := c.Params.Get("fullname")
	avatar := c.Params.Files["avatar"]
	email := c.Params.Get("email")
	currpassword := c.Params.Get("currpassword")
	password := c.Params.Get("password")
	confirmPassword := c.Params.Get("confirmpassword")
	mobile := c.Params.Get("mobile")
	address := c.Params.Get("address")

	hus := models.Hus{}
	husAjax := models.HusAjax{}

	//Validate
	if !c.Validation.Required(fullname).Ok {
		husAjax.SetMessage("Vui lòng nhập Họ và tên")
		return c.RenderJSON(husAjax.OutData(false))
	}
	fileTypes := []string{"image/png", "image/jpg", "image/jpeg", "gif"}
	if len(avatar) != 0 && !hus.ValidateFiles(avatar, fileTypes) {
		husAjax.SetMessage("Avatar sai định dạng.")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Required(email).Ok {
		husAjax.SetMessage("Vui lòng nhập Email")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Email(email).Ok {
		husAjax.SetMessage("Vui lòng nhập đúng chuẩn Email")
		return c.RenderJSON(husAjax.OutData(false))
	}

	if id == 0 && !c.Validation.Required(password).Ok {
		husAjax.SetMessage("Vui lòng nhập Mật khẩu")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if id != 0 && !c.Validation.Required(currpassword).Ok {
		husAjax.SetMessage("Vui lòng nhập  Mật khẩu")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if len(password) > 0 && !c.Validation.MinSize(password, 8).Ok {
		husAjax.SetMessage("Mật khẩu có ít nhất 8 ký tự")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if password != confirmPassword {
		husAjax.SetMessage("Nhập lại mật khẩu không trùng khớp mật khẩu đã nhập")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Required(mobile).Ok {
		husAjax.SetMessage("Vui lòng nhập Số điện thoại")
		return c.RenderJSON(husAjax.OutData(false))
	}
	//All data validated
	repo := dataservice.NewUserRepo()
	user := repo.GetByEmail(email)

	token := ""
	if id == 0 {// Register
		responseCaptcha := c.Params.Form.Get("g-recaptcha-response")
		secret, _ := revel.Config.String("CAPTCHA_SECRET_KEY")

		husCaptcha := models.HusCaptcha{Secret: secret}

		isValid := husCaptcha.Verify(responseCaptcha)

		if !isValid {
			husAjax.SetMessage("Vui lòng xác thực Captcha thành công.")
			return c.RenderJSON(husAjax.OutData(false))
		}

		//Check user if exist
		if user.Id > 0 && user.Status == 1 {
			husAjax.SetMessage("Email đã tồn tại trong hệ thống.")
			return c.RenderJSON(husAjax.OutData(false))
		}

		//Generate Token
		token = hus.RandomString(16)

		husMail := models.HusMail{"noreply@myphamtayninh.com", "!@#123^%$"}
		//Supply information to send email
		receivers := []string{email}
		subject := "Kích hoạt tài khoản Mỹ Phẩm Tây Ninh"
		bodyMessage := "Chào mừng quý khách đến với Mỹ Phẩm Tây Ninh.<br/>Vui lòng click vào <a href=\"" + hus.BaseUrl() + "/user/" + email + "/" + token + "\">đây</a> để kích hoạt tài khoản của bạn."
		//Send email
		err := husMail.SendMail(receivers, subject, bodyMessage)

		if err != nil {
			husAjax.SetMessage("Gửi mail xác thực không thành công.")
			return c.RenderJSON(husAjax.OutData(false))
		}

		c.Flash.Success("Vui lòng kích hoạt tài khoản từ địa chỉ email đã đăng ký.")
	}

	//Encrypt md5 password
	h1 := md5.New()
	io.WriteString(h1, password)
	encryptPassword := hex.EncodeToString(h1.Sum(nil))
	//Encrypt md5 currpassword
	h2 := md5.New()
	io.WriteString(h2, currpassword)
	encryptCurrPassword := hex.EncodeToString(h2.Sum(nil))

	var userEmail string
	if user.Id == 0 {
		user := models.User{FullName: fullname, Email: email, Password: encryptPassword, Mobile: mobile, Address: address, Token: token}
		repo.Create(&user)
		userEmail = user.Email
	} else {
		//Check currPassWord
		if user.Password != encryptCurrPassword {
			husAjax.SetMessage("Sai mật khẩu hiện tại.")
			return c.RenderJSON(husAjax.OutData(false))
		}

		user.FullName = fullname
		if len(password) > 0 {
			user.Password = encryptPassword
		}
		user.Mobile = mobile
		user.Address = address
		user.Token = token
		repo.Update(user)
		userEmail = user.Email
		if id > 0 {
			//Update session in case user logged in
			jsonUser, _ := json.Marshal(user)
			c.Session["loggedUser"] = string(jsonUser)
			c.Flash.Success("Cập nhật thông tin tài khoản của bạn thành công.")
		}
	}

	if len(avatar) > 0 {
		user = repo.GetByEmail(userEmail)
		//Delete old avatar
		hus.DeleteDirFile(user.Avatar)

		//Upload and resize avatar
		uploadDir := "/public/images/avatars"
		avatarFile := models.HusFile{uploadDir, avatar}
		avatarName := "avatar_" + strconv.FormatUint(uint64(user.Id), 10) + path.Ext(avatarFile.Files[0].Filename)
		uploadAvatar := avatarFile.UploadFile(strings.ToLower(avatarName));
		if uploadAvatar {
			avatarFile.ThumbnailImage(avatarName, 80, 80)
			user.Avatar = uploadDir+"/"+avatarName;
			repo.Update(user)
			if id > 0 {
				//Update session in case user logged in
				jsonUser, _ := json.Marshal(user)
				c.Session["loggedUser"] = string(jsonUser)
			}
		}
	}

	return c.RenderJSON(husAjax.OutData(true))
}

func (c User) ForgotPasswordForm() revel.Result {
	husAjax := models.HusAjax{}

	var data interface{}
	content := husAjax.Fetch("User/ForgotPasswordForm.html", data)
	husAjax.SetHTML("common_dialog", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c User) ForgotPasswordProcess() revel.Result {
	husAjax := models.HusAjax{}

	email := c.Params.Get("email")
	if !c.Validation.Required(email).Ok {
		husAjax.SetMessage("Vui lòng nhập Email")
		return c.RenderJSON(husAjax.OutData(false))
	}
	if !c.Validation.Email(email).Ok {
		husAjax.SetMessage("Vui lòng nhập đúng chuẩn Email")
		return c.RenderJSON(husAjax.OutData(false))
	}

	repo := dataservice.NewUserRepo()
	//Check user if exist
	user := repo.GetByEmail(email)
	if user.Id == 0 || user.Status != 1 {
		husAjax.SetMessage("Email không tồn tại trong hệ thống.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	//Generate token
	hus := models.Hus{}
	token := hus.RandomString(16)

	husMail := models.HusMail{"noreply@myphamtayninh.com", "!@#123^%$"}

	//Supply information to send email
	receivers:= []string{email}
	subject := "Mật khẩu mới tài khoản Mỹ Phẩm Tây Ninh"
	bodyMessage := "Quý khách quên mật khẩu tài khoản Mỹ Phẩm Tây Ninh?<br/>Vui lòng click vào <a href=\""+hus.BaseUrl()+"/forgotpassword/"+email+"/"+token+"\">đây</a> để thay đổi mật khẩu của bạn."
	//Send email
	err := husMail.SendMail(receivers, subject, bodyMessage)

	if err != nil {
		husAjax.SetMessage("Gửi mail không thành công.")
		return c.RenderJSON(husAjax.OutData(false))
	}

	user.Token = token
	repo.Update(user)

	c.Flash.Success("Vui lòng kiểm tra email của bạn để đổi mật khẩu mới.")

	return c.RenderJSON(husAjax.OutData(true))
}

func (c User) ForgotPasswordChangeForm() revel.Result {
	email := c.Params.Get("email")
	token := c.Params.Get("token")

	//Check token
	repo := dataservice.NewUserRepo()
	user := repo.GetByEmail(email)

	if user.Id == 0 {
		c.Flash.Error("Tài khoản không tồn tại trong hệ thống.")
		return  c.Redirect(routes.App.Index())
	}

	if user.Token != token {
		c.Flash.Error("Hết hạn token.")
		return  c.Redirect(routes.App.Index())
	}

	return c.Render(email, token)
}

func (c User) ForgotPasswordUpdate() revel.Result {
	email := c.Params.Get("email")
	password := c.Params.Get("password")
	confirmPassword := c.Params.Get("confirmpassword")
	token := c.Params.Get("token")

	repo := dataservice.NewUserRepo()
	user := repo.GetByEmail(email)

	//Validate
	c.Validation.Required(email).Message("Vui lòng nhập Email")
	c.Validation.Email(email).Message("Vui lòng nhập đúng chuẩn Email")
	c.Validation.Required(password).Message("Vui lòng nhập Password")
	c.Validation.MinSize(password, 8).Message("Password có ít nhất 8 ký tự")
	c.Validation.Required(password == confirmPassword).Message("Nhập lại mật khẩu không trùng khớp mật khẩu đã nhập.")
	c.Validation.Required(user.Id > 0).Message("Tài khoản không tồn tại trong hệ thống.")

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		c.ViewArgs["email"] = email
		c.ViewArgs["token"] = token
		return c.RenderTemplate("User/ForgotPasswordChangeForm.html")
	}
	//All data validated

	//Check token
	if user.Token != token {
		c.Flash.Error("Hết hạn token.")
		return  c.Redirect(routes.App.Index())
	}

	//Encrypt md5 password
	h := md5.New()
	io.WriteString(h, password)
	encryptPassword := hex.EncodeToString(h.Sum(nil))

	user.Password = encryptPassword
	user.Token = ""
	repo.Update(user)

	c.Flash.Success("Thay đổi mật khẩu mới thành công.")
	return c.Redirect(routes.User.Login())
}
