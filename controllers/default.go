package controllers

import (
	"StockMe/models"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// API for Client
type API struct {
	beego.Controller
}

type ResponseObject struct {
	Code           int         `json:"code"`
	Message        string      `json:"message"`
	ResponseObject interface{} `json:"responseObject"`
}

type ResponseObjectWithCode struct {
	Code           int64       `json:"code"`
	Message        string      `json:"message"`
	ResponseObject interface{} `json:"responseObject"`
}

var (
// defAdminPassword = "password"
)

func Init() {
	logs.Debug("init db")
	initBaseServer()
}

func (api *API) BaseURL() string {
	var baseUrl string = api.Ctx.Input.Site() + fmt.Sprintf(":%d", api.Ctx.Input.Port())
	if api.Ctx.Input.Header("X-Forwarded-Host") != "" {
		baseUrl = api.Ctx.Input.Scheme() + "://" + api.Ctx.Input.Header("X-Forwarded-Host")
	}
	return baseUrl
}

func (api *API) ResponseJSON(results interface{}, code int, msg string) {
	if results == nil {
		results = struct{}{}
	}
	response := &ResponseObject{
		Code:           code,
		Message:        msg,
		ResponseObject: results,
	}
	api.Data["json"] = response
	api.Ctx.Output.SetStatus(code)
	api.ServeJSON()
	return
}

func (api *API) ResponseJSONWithCode(results interface{}, statusCode int, code int64, msg string) {
	if results == nil {
		results = struct{}{}
	}
	response := &ResponseObjectWithCode{
		Code:           code,
		Message:        msg,
		ResponseObject: results,
	}
	api.Data["json"] = response
	api.Ctx.Output.SetStatus(statusCode)
	api.ServeJSON()
	return
}

func (c *API) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func initBaseServer() {
	c := "usd"
	existC, _ := models.GetCurrencyByTitle(c)

	if existC == nil {
		cId, err := models.AddCurrency(&models.Currency{
			Title:   c,
			Symbol:  "$",
			Divisor: 0,
		})
		if err != nil {
			logs.Error("err init Currency:", c, "err", err)
		}
		existC, _ = models.GetCurrencyById(cId)
	}

	appName := "StockMe"
	exist, _ := models.GetUserByUsername(appName)
	if exist == nil {
		_, err := models.AddUser(&models.User{
			Id:            0,
			Currency:      existC,
			Username:      appName,
			Password:      "1234",
			FirebaseToken: "",
			FirebaseUuid:  "",
			Activate:      true,
		})
		if err != nil {
			logs.Error("err init position:", appName, "err", err)
		}
	}

}
