package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Assets struct {
	Id      int64     `orm:"auto"`
	User    *User     `orm:"null;rel(fk)"`
	Stock   *Stock    `orm:"null;rel(fk)"`
	Created time.Time `orm:"auto_now_add;type(datetime)" json:"created"`
	Updated time.Time `orm:"auto_now;type(datetime)" json:"updated"`
}
type ImportAssets struct {
	PurchaseDate  string  `json:"purchase_date"`
	Symbol        string  `json:"symbol"`
	Title         string  `json:"title"`
	PurchasePrice float64 `json:"purchasePrice"`
	Investment    float64 `json:"investment"`
	Amount        float64 `json:"amount"`
	ThaiBaht      float64 `json:"thaiBaht"`
}

func init() {
	orm.RegisterModel(new(Assets))
}

// AddAssets insert a new Assets into database and returns
// last inserted Id on success.
func AddAssets(m *Assets) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetAssetsById(id int64) (v *Assets, err error) {
	o := orm.NewOrm()
	v = &Assets{Id: id}
	if err = o.QueryTable(new(Assets)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateAssetsById(o orm.Ormer, m *Assets) (err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	v := Assets{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteAssets(id int64) (err error) {
	o := orm.NewOrm()
	v := Assets{Id: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Assets{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAssetsList(limit, offset int64) (v []*Assets, total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Assets))

	total, _ = qs.RelatedSel().Count()
	_, err = qs.RelatedSel().Limit(limit, offset).All(&v)
	if err == nil {
		return v, total, nil
	}
	return nil, total, err
}

func GetAssetsByStock(stockId int64) (v *Assets, err error) {
	o := orm.NewOrm()
	v = &Assets{}
	if err = o.QueryTable(new(Assets)).Filter("Stock", stockId).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}
