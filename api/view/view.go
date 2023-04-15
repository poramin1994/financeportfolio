package view

import (
	v1 "StockMe/api"
	util "StockMe/api/util"
	"StockMe/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type View struct {
	v1.API
}
type Response struct {
	GlobalQuote GlobalQuote `json:"Global Quote"`
}
type GlobalQuote struct {
	Symbol           string `json:"01. symbol"`
	Open             string `json:"02. open"`
	High             string `json:"03. high"`
	Low              string `json:"04. low"`
	Price            string `json:"05. price"`
	Volume           string `json:"06. volume"`
	LatestTradingDay string `json:"07. latest trading day"`
	PreviousClose    string `json:"08. previous close"`
	Change           string `json:"09. change"`
	ChangePercent    string `json:"10. change percent"`
}

func (c *View) GetHomePage() {
	assetsList := make([]map[string]interface{}, 0)

	assets, _, err := models.GetAssetsList(0, -1)
	if err != nil {
		c.ResponseJSONWithCode(nil, 500, 50000, "error getting assets")
		return
	}
	for _, v := range assets {
		resp, err := GetPrice5Min(v.Stock.Symbol)
		if err != nil {
			c.ResponseJSONWithCode(nil, 500, 50001, "error get price")
			return
		}
		assetsList = append(assetsList, genAssetsListItem(v, &resp))
	}

	c.Data["Website"] = "My Strock"
	c.Data["Email"] = "poramin.m1994@gmail.com"
	c.Data["Done"] = false
	c.Data["assetsList"] = assetsList
	c.TplName = "index.html"
}

func (c *View) GetAssetsDetailPage() {
	symbol, _ := c.GetInt64("symbol", 0)
	if symbol == 0 {
		c.ResponseJSONWithCode(nil, 400, 40000, "error bad request")
		return
	}
	assets, _ := models.GetAssetsById(symbol)

	assetsList := make([]map[string]interface{}, 0)
	assetsDetail, _, err := models.GetAssetsDetailByAssetIdList(symbol)
	if err != nil {
		c.ResponseJSONWithCode(nil, 500, 50000, "error getting assets")
		return
	}
	for _, v := range assetsDetail {
		resp, err := GetPrice5Min(assets.Stock.Symbol)
		if err != nil {
			c.ResponseJSONWithCode(nil, 500, 50001, "error get price")
			return
		}
		assetsList = append(assetsList, genAssetsDetailsItem(v, &resp))
	}

	c.Data["Website"] = "My Strock"
	c.Data["Email"] = "poramin.m1994@gmail.com"
	c.Data["Done"] = false
	c.Data["assetsList"] = assetsList
	c.TplName = "detailAsset.html"
}

func GetPrice5Min(symbol string) (Response, error) {
	var respond Response
	var err error
	txn := util.FormatDateUnix(time.Now())
	var url = fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, txn)
	_, bodyRespond, _ := util.CallPostHttpWithJsonRespond(url, nil, nil, "GET")
	bodyRespondByte, err := json.Marshal(bodyRespond)
	if err != nil {
		return respond, err
	}
	err = json.Unmarshal(bodyRespondByte, &respond)
	if err != nil {
		return respond, err
	}
	return respond, err
}

func genAssetsListItem(s *models.Assets, response *Response) map[string]interface{} {
	price, _ := strconv.ParseFloat(response.GlobalQuote.Price, 64)
	// change := response.GlobalQuote.Change
	change, _ := strconv.ParseFloat(response.GlobalQuote.Change, 64)
	changePercent := response.GlobalQuote.ChangePercent
	res1 := strings.Split(response.GlobalQuote.Change, "")
	sId := strconv.FormatInt(s.Id, 10)
	urlDetails := "/web/assets/detail?symbol=" + sId
	class := "positive"
	if len(res1) > 0 {
		if res1[0] == "-" {
			class = "negative"
		}
	}
	taskDetail := map[string]interface{}{
		"title":         s.Stock.Title,
		"symbol":        s.Stock.Symbol,
		"price":         price,
		"change":        change,
		"changePercent": changePercent,
		"class":         class,
		"urlDetails":    urlDetails,
	}

	return taskDetail
}

func genAssetsDetailsItem(s *models.AssetsDetail, response *Response) map[string]interface{} {
	price, _ := strconv.ParseFloat(response.GlobalQuote.Price, 64)

	finalValue := (price - s.PurchasePrice) * s.Amount
	stringFinalValue := fmt.Sprintf("%.2f", finalValue)
	class := "positive"
	symbol := "+"
	if finalValue < 0 {
		class = "negative"
		symbol = ""
	}
	taskDetail := map[string]interface{}{
		"class":         class,
		"purchaseDate":  s.PurchaseDate.Format("2006-01-02"),
		"purchasePrice": s.PurchasePrice,
		"initialValue":  s.InitialValue,
		"amount":        s.Amount,
		"thaiBaht":      s.ThaiBaht,
		"finalValue":    symbol + stringFinalValue,
	}

	return taskDetail
}
