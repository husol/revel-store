package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/dataservice"
	"crypto/md5"
	"io"
	"encoding/hex"
	"husol.org/mypham/app/routes"
)

type AdminUser struct {
	*revel.Controller
}

func (c AdminUser) Index() revel.Result {
	//Count new comments
	repoComment := dataservice.NewCommentRepo()
	newCommentNum := repoComment.CountNewComments()
	c.ViewArgs["newCommentNum"] = newCommentNum

	repo := dataservice.NewUserRepo()
	var users []models.User
	users, _ = repo.GetAll()
	c.ViewArgs["users"] = users

	return c.Render()
}

func (c AdminUser) LoadForm() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repo := dataservice.NewUserRepo()
	var data interface{}

	data = models.User{}
	if id > 0 {
		user := repo.GetById(id)
		data = user
	}

	husAjax := models.HusAjax{}
	content := husAjax.Fetch("AdminUser/Form.html", data)
	husAjax.SetHTML("common_dialog", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c AdminUser) Update() revel.Result {
	var id, role, status int
	c.Params.Bind(&id, "id_record")
	fullName := c.Params.Get("fullname")
	email := c.Params.Get("email")
	password := c.Params.Get("password")
	mobile := c.Params.Get("mobile")
	address := c.Params.Get("address")
	c.Params.Bind(&role, "role")
	c.Params.Bind(&status, "status")

	//Encrypt md5 password
	h := md5.New()
	io.WriteString(h, password)
	encryptPassword := hex.EncodeToString(h.Sum(nil))

	repo := dataservice.NewUserRepo()

	if id > 0 {
		user := repo.GetById(id)
		user.FullName = fullName
		user.Email = email
		if (password != "") {
			user.Password = encryptPassword
		}
		user.Mobile = mobile
		user.Address = address
		user.Role = role
		user.Status = status

		repo.Update(user);
		c.Flash.Success("Cập nhật User thành công.")
	} else {
		user := models.User{FullName: fullName, Email: email, Password: encryptPassword, Mobile: mobile, Address: address, Role: role, Status: status}

		repo.Create(&user)
		c.Flash.Success("Thêm User thành công.")
	}

	return c.Redirect("/admin/users")
}

func (c AdminUser) Delete() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	hus := models.Hus{}
	repo := dataservice.NewUserRepo()
	user := repo.GetById(id)

	if (user.Id > 0) {
		//Delete avatar
		hus.DeleteDirFile(user.Avatar)
		repo.Delete(user)
		c.Flash.Success("Xóa User thành công.")
	} else {
		c.Flash.Error("Không thể thực hiện xóa User")
	}
	return c.Redirect(routes.AdminUser.Index())
}

func (c AdminUser) CheckExistedEmail() revel.Result {
	var id uint
	c.Params.Bind(&id, "id_record")
	email := c.Params.Get("email")

	husAjax := models.HusAjax{}
	repo := dataservice.NewUserRepo();
	user := repo.GetByEmail(email)
	result := user.Id > 0
	if id > 0 && user.Id == id {
		result = false
	}

	return c.RenderJSON(husAjax.OutData(result))
}