package manager

import (
	v1 "StockMe/api"
	util "StockMe/api/util"
	"StockMe/models"
	"encoding/json"
	"time"

	"github.com/beego/beego/logs"
)

type Assets struct {
	v1.API
}

func (this *Assets) AddAssets() {
	result := map[string]interface{}{}
	var data models.ImportAssets
	now := time.Now()
	user, _ := models.GetUserById(1)

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.Error(" ImportGenTask : json.Unmarshal ## err = ", err)
		this.ResponseJSONWithCode(result, 500, 50000, err.Error())
		return
	}

	stockObject, _ := models.GetStockBySymbol(data.Symbol)
	if stockObject == nil {
		sId, err := models.AddStock(&models.Stock{
			Title:   data.Title,
			Symbol:  data.Symbol,
			Created: now,
			Updated: now,
		})
		if err != nil {
			this.ResponseJSONWithCode(result, 500, 50001, err.Error())
			return
		}
		stockObject, _ = models.GetStockById(sId)
	}

	assetObject, _ := models.GetAssetsByStock(stockObject.Id)
	if assetObject == nil {
		aId, err := models.AddAssets(&models.Assets{
			User:    user,
			Stock:   stockObject,
			Created: now,
			Updated: now,
		})
		if err != nil {
			this.ResponseJSONWithCode(result, 500, 50002, err.Error())
			return
		}
		assetObject, _ = models.GetAssetsById(aId)
	}

	_, err = models.AddAssetsDetail(&models.AssetsDetail{
		Assets:        assetObject,
		Amount:        data.Amount,
		InitialValue:  data.Investment,
		PurchasePrice: data.PurchasePrice,
		ThaiBaht:      data.ThaiBaht,
		PurchaseDate:  util.ToDateTime(data.PurchaseDate),
		Created:       now,
		Updated:       now,
	})
	if err != nil {
		this.ResponseJSONWithCode(result, 500, 50003, err.Error())
		return
	}

	this.ResponseJSONWithCode(result, 200, 200, "Successfully")
	return
}
