package controllers

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/routes"
	"math"
)

type AdminComment struct {
	*revel.Controller
}

const COMMENT_PAGESIZE int = 10

func (c AdminComment) Index() revel.Result {
	//Count new comments
	repoComment := dataservice.NewCommentRepo()
	newCommentNum := repoComment.CountNewComments()
	c.ViewArgs["newCommentNum"] = newCommentNum

	//Calculate total page
	comments := repoComment.Search("", -1, -1)
	totalPage := math.Ceil(float64(len(comments)) / float64(COMMENT_PAGESIZE))
	c.ViewArgs["totalPage"] = int(totalPage)

	return c.Render()
}

func (c AdminComment) Paging() revel.Result {
	var page int
	c.Params.Bind(&page, "page")

	repoComment := dataservice.NewCommentRepo()
	comments := repoComment.Search("", page, COMMENT_PAGESIZE)

	husAjax := models.HusAjax{}
	content := husAjax.Fetch("AdminComment/list.html", comments)
	husAjax.SetHTML("tableComment", content)

	return c.RenderJSON(husAjax.OutData(true))
}

func (c AdminComment) Approve() revel.Result {
	var id int
	c.Params.Bind(&id, "id")
	repo := dataservice.NewCommentRepo()

	comment := repo.GetById(id)
	if comment.Id > 0 {
		comment.Status = 1

		repo.Update(comment)

		c.Flash.Success("Đã duyệt Bình luận thành công.")
	} else {
		c.Flash.Error("Bình luận không tồn tại.")
	}

	return c.Redirect(routes.AdminComment.Index())
}

func (c AdminComment) Delete() revel.Result {
	var id int
	c.Params.Bind(&id, "id")

	repo := dataservice.NewCommentRepo()
	comment := repo.GetById(id)

	if comment.Id > 0 {
		repo.Delete(comment)
		c.Flash.Success("Xóa Bình luận thành công.")
	} else {
		c.Flash.Error("Không thể thực hiện xóa Bình luận.")
	}
	return c.Redirect(routes.AdminComment.Index())
}