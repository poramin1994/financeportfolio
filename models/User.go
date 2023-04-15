package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id            int64     `orm:"auto"`
	Currency      *Currency `orm:"null;rel(fk)"`
	Username      string    `orm:"null"`
	Password      string    `orm:"null"`
	FirebaseToken string    `orm:"null;type(text)"`
	FirebaseUuid  string    `orm:"null"`
	Activate      bool      `orm:"default(0)"`
	Delete        bool      `orm:"default(0)"`
	Deleted       time.Time `orm:"null"`
	DeletedBy     int64     `orm:"null"`

	Created time.Time `orm:"auto_now_add;type(datetime)" json:"created"`
	Updated time.Time `orm:"auto_now;type(datetime)" json:"updated"`
}

func init() {
	orm.RegisterModel(new(User))
}

func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetUserById(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.QueryTable(new(User)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateUserById(o orm.Ormer, m *User) (err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	v := User{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetUserByUsername(username string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{}
	if err = o.QueryTable(new(User)).Filter("Username", username).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}
