package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// create a new struct for new dadatbase table with Finance Stock
type Stock struct {
	Id      int64     `orm:"auto"`
	Title   string    `orm:"null"`
	Symbol  string    `orm:"null"`
	Created time.Time `orm:"auto_now_add;type(datetime)" json:"created"`
	Updated time.Time `orm:"auto_now;type(datetime)" json:"updated"`
}

func init() {
	orm.RegisterModel(new(Stock))
}

func AddStock(m *Stock) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetStockById(id int64) (v *Stock, err error) {
	o := orm.NewOrm()
	v = &Stock{Id: id}
	if err = o.QueryTable(new(Stock)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateStockById(o orm.Ormer, m *Stock) (err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	v := Stock{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteStock(id int64) (err error) {
	o := orm.NewOrm()
	v := Stock{Id: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Stock{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetStockBySymbol(symbol string) (v *Stock, err error) {
	o := orm.NewOrm()
	v = &Stock{}
	if err = o.QueryTable(new(Stock)).Filter("Symbol", symbol).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}
