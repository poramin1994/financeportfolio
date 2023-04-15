package routers

import (
	v1api "StockMe/api/v1/api"
	v1view "StockMe/api/view"
	"html/template"
	"net/http"

	"StockMe/controllers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.Options("/*", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
		ctx.Output.Body([]byte("OK"))
		return
	})
	beego.ErrorHandler("404", page_not_found)
	beego.Router("/", &controllers.API{})

	managev1 := beego.NewNamespace("/v1/manage/api",
		beego.NSBefore(FilterDebug),
		beego.NSNamespace("/assets",
			beego.NSRouter("/add", &v1api.Assets{}, "Post:AddAssets"),
		),
	)
	beego.AddNamespace(managev1)

	view := beego.NewNamespace("/web",
		beego.NSBefore(FilterDebug),
		beego.NSNamespace("/home",
			beego.NSRouter("/", &v1view.View{}, "get:GetHomePage"),
		),
		beego.NSNamespace("/assets",
			beego.NSRouter("/detail", &v1view.View{}, "get:GetAssetsDetailPage"),
		),
	)
	beego.AddNamespace(view)

}

var FilterDebug = func(ctx *context.Context) {
	logs.Debug("--------------------------------------------")
	logs.Debug("RequestURI", ctx.Request.RequestURI)
	logs.Debug("method", ctx.Request.Method)
	logs.Debug("content-type:", ctx.Request.Header.Values("Content-Type"))
	logs.Debug("params:", ctx.Input.Params())
	logs.Debug("body:", string(ctx.Input.RequestBody))
	logs.Debug("form:", ctx.Request.Form)
	logs.Debug("postform:", ctx.Request.PostForm)
	logs.Debug("=== Request Headers ===")
	for name, value := range ctx.Request.Header {
		logs.Debug(name, ":", value)
	}
	logs.Debug("=== END Request Headers ===")
	return
}

func page_not_found(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/404.html")
	data := make(map[string]interface{})
	data["content"] = "page not found"
	t.Execute(rw, data)
}
