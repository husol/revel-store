package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/routes"
	"strconv"
	"path"
	"strings"
)

type AdminInformation struct {
	*revel.Controller
}

func (c AdminInformation) Index() revel.Result {
	//Count new comments
	repoComment := dataservice.NewCommentRepo()
	newCommentNum := repoComment.CountNewComments()
	c.ViewArgs["newCommentNum"] = newCommentNum

	repo := dataservice.NewInformationRepo()
	information := repo.GetAll()
	c.ViewArgs["information"] = information

	return c.Render()
}

func (c AdminInformation) LoadForm() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repo := dataservice.NewInformationRepo()
	var data interface{}

	data = models.Information{}
	if id > 0 {
		information := repo.GetById(id)
		data = information
	}

	husAjax := models.HusAjax{}
	content := husAjax.Fetch("AdminInformation/Form.html", data)
	husAjax.SetHTML("common_dialog", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c AdminInformation) Update() revel.Result {
	var id, status int
	c.Params.Bind(&id, "id_record")
	title := c.Params.Get("title")
	cover := c.Params.Files["cover"]
	description := c.Params.Get("description")
	c.Params.Bind(&status,"status")
	repo := dataservice.NewInformationRepo()

	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	if id > 0 {
		information := repo.GetById(id)
		information.Title = title
		information.Description = description
		information.Status = status

		if len(cover) > 0 {
			//Delete old cover
			hus.DeleteDirFile(information.Cover)

			//Upload and resize cover
			uploadDir := "/public/images/information"
			coverFile := models.HusFile{uploadDir, cover}
			coverName := "info_" + strconv.FormatUint(uint64(information.Id), 10) + path.Ext(coverFile.Files[0].Filename)
			uploadAvatar := coverFile.UploadFile(strings.ToLower(coverName));
			if uploadAvatar {
				coverFile.ThumbnailImage(coverName, 500, 315)
				information.Cover = uploadDir+"/"+coverName;
			}
		}
		repo.Update(information)

		c.Flash.Success("Cập nhật Thông tin thành công.")
	} else {
		information := models.Information{IdUser:loggedUser.Id, Title: title, Description: description, Status: status}
		repo.Create(&information)

		if len(cover) > 0 {
			//Upload and resize cover
			uploadDir := "/public/images/information"
			coverFile := models.HusFile{uploadDir, cover}
			coverName := "info_" + strconv.FormatUint(uint64(information.Id), 10) + path.Ext(coverFile.Files[0].Filename)
			uploadAvatar := coverFile.UploadFile(strings.ToLower(coverName));
			if uploadAvatar {
				coverFile.ThumbnailImage(coverName, 500, 315)
				information.Cover = uploadDir + "/" + coverName;
				repo.Update(&information)
			}
		}
		c.Flash.Success("Thêm Thông tin thành công.")
	}

	return c.Redirect(routes.AdminInformation.Index())
}

func (c AdminInformation) Delete() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	hus := models.Hus{}
	repo := dataservice.NewInformationRepo()
	information := repo.GetById(id)

	if information.Id > 0 {
		//Delete cover
		hus.DeleteDirFile(information.Cover)
		repo.Delete(information)
		c.Flash.Success("Xóa Thông tin thành công.")
	} else {
		c.Flash.Error("Không thể thực hiện xóa Thông tin.")
	}
	return c.Redirect(routes.AdminInformation.Index())
}