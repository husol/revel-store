package app

import (
	"github.com/revel/revel"
	"husol.org/mypham/app/dataservice"
	"husol.org/mypham/app/models"
	"husol.org/mypham/app/controllers"
)

func prependRoute(c *revel.Controller) revel.Result {
	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	c.ViewArgs["_loggedUser"] = loggedUser
	c.ViewArgs["_curr_controller"] = c.Name
	c.ViewArgs["_curr_method"] = c.MethodName
	c.ViewArgs["_curr_route"] = c.Action

	return nil
}

func checkUser(c *revel.Controller) revel.Result {
	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	if loggedUser.Id == 0 {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			return c.RenderJSON("expired_session")
		}
		c.Flash.Error("Vui lòng đăng nhập.")
		return c.Redirect("/login")
	}
	return nil
}

func checkAdminUser(c *revel.Controller) revel.Result {
	hus := models.Hus{}
	loggedUser := models.User{}
	hus.DecodeObjSession(c.Session["loggedUser"], &loggedUser)

	if loggedUser.Id == 0 || loggedUser.Role != 0 {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			return c.RenderJSON("expired_session")
		}
		c.Flash.Error("Vui lòng đăng nhập với quyền Admin.")
		return c.Redirect("/login")
	}

	return nil
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(InitDB) // invoke InitDB function before
	//Custom template functions
	husHtml := models.HusHtml{}
	revel.TemplateFuncs["html"] = husHtml.Html
	revel.TemplateFuncs["mod"] = husHtml.Mod
	revel.TemplateFuncs["formatFValue"] = husHtml.FormatFValue
	revel.TemplateFuncs["add"] = husHtml.Add
	revel.TemplateFuncs["subtract"] = husHtml.Subtract
	revel.TemplateFuncs["formatCurrTime"] = husHtml.FormatCurrTime
	revel.TemplateFuncs["inArray"] = husHtml.InArray

	revel.InterceptMethod(prependRoute, revel.AFTER)

	revel.InterceptFunc(checkAdminUser, revel.BEFORE, &controllers.AdminHome{})
	revel.InterceptFunc(checkAdminUser, revel.BEFORE, &controllers.AdminTransaction{})
	revel.InterceptFunc(checkAdminUser, revel.BEFORE, &controllers.AdminUser{})
	revel.InterceptFunc(checkAdminUser, revel.BEFORE, &controllers.AdminCategory{})
	revel.InterceptFunc(checkAdminUser, revel.BEFORE, &controllers.AdminProduct{})
	revel.InterceptFunc(checkAdminUser, revel.BEFORE, &controllers.AdminInformation{})
	revel.InterceptFunc(checkAdminUser, revel.BEFORE, &controllers.AdminComment{})

}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func InitDB() {
	driver := revel.Config.StringDefault("db.driver", "mysql")
	spec := revel.Config.StringDefault("db.spec", "")
	dataservice.InitDb(driver, spec)
}
