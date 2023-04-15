package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type AssetsDetail struct {
	Id            int64     `orm:"auto"`
	Assets        *Assets   `orm:"null;rel(fk)"`
	Amount        float64   `orm:"null"`
	InitialValue  float64   `orm:"null"`
	PurchasePrice float64   `orm:"null"`
	ThaiBaht      float64   `orm:"null"`
	PurchaseDate  time.Time `orm:"auto_now_add;type(datetime)"`
	Created       time.Time `orm:"auto_now_add;type(datetime)" json:"created"`
	Updated       time.Time `orm:"auto_now;type(datetime)" json:"updated"`
}

func init() {
	orm.RegisterModel(new(AssetsDetail))
}

func AddAssetsDetail(m *AssetsDetail) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetAssetsDetailById(id int64) (v *AssetsDetail, err error) {
	o := orm.NewOrm()
	v = &AssetsDetail{Id: id}
	if err = o.QueryTable(new(AssetsDetail)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateAssetsDetailById(o orm.Ormer, m *AssetsDetail) (err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	v := AssetsDetail{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteAssetsDetail(id int64) (err error) {
	o := orm.NewOrm()
	v := AssetsDetail{Id: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AssetsDetail{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAssetsDetailByAssetIdList(taskId int64) (v []*AssetsDetail, total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AssetsDetail))

	total, _ = qs.Filter("Assets__Id", taskId).RelatedSel().Count()
	_, err = qs.Filter("Assets__Id", taskId).RelatedSel().All(&v)
	if err == nil {
		return v, total, nil
	}
	return nil, total, err
}
