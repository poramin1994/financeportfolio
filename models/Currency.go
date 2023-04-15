package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Currency struct {
	Id      int64     `orm:"auto"`
	Title   string    `orm:"null"`
	Symbol  string    `orm:"null"`
	Divisor float64   `orm:"null"`
	Created time.Time `orm:"auto_now_add;type(datetime)" json:"created"`
	Updated time.Time `orm:"auto_now;type(datetime)" json:"updated"`
}

func init() {
	orm.RegisterModel(new(Currency))
}

func AddCurrency(m *Currency) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetCurrencyById(id int64) (v *Currency, err error) {
	o := orm.NewOrm()
	v = &Currency{Id: id}
	if err = o.QueryTable(new(Currency)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateCurrencyById(o orm.Ormer, m *Currency) (err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	v := Currency{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteCurrency(id int64) (err error) {
	o := orm.NewOrm()
	v := Currency{Id: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Currency{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetCurrencyByTitle(title string) (v *Currency, err error) {
	o := orm.NewOrm()
	v = &Currency{}
	if err = o.QueryTable(new(Currency)).Filter("Title", title).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}
