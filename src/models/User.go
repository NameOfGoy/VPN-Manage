package models

import (
	db "VPN-Manage/src/database"
	e "VPN-Manage/src/util/myerror"
	"log"
	"fmt"
)

type User struct {
	Uuid     string `xorm:"not null pk VARCHAR(32)"`
	Loginid  string `xorm:"not null VARCHAR(20)"`
	Password string `xorm:"not null VARCHAR(20)"`
	Username string `xorm:"VARCHAR(255)"`
	Isbindwx string `xorm:"VARCHAR(255)"`
}

func (u *User)AddUser() (err error) {
	_,err = db.Engine.Insert(u)
	if err != nil {
		errmsg := fmt.Sprintf("%s", err)
		err = &e.Myerror{
			Code : e.ERROR_INSERT_TABLE_CODE,
			Message : "User表插入错误: " + errmsg,
		}
		return
	}
	return
}

func (u *User)GetUser() (user User,err error) {
	ret,err := db.Engine.Alias("user").Where("user.loginid = ?", u.Loginid).Get(&user)
	if err != nil{
		errmsg := fmt.Sprintf("%s", err)
		err = &e.Myerror{
			Code : e.ERROR_SELECT_TABLE_CODE,
			Message : "user表查询错误:" + errmsg,
		}
		return
	}
	if !ret {
		err = &e.Myerror{
			Code : e.ERROR_USER_NOTFOUND_CODE,
			Message : e.ERROR_USER_NOTFOUND_MESSAGE,
		}
	}
	return
}

func (u *User)UpdateUser()(err error) {
	_, err = db.Engine.Exec("UPDATE USER SET PASSWORD = ? WHERE LOGINID = ?", u.Password, u.Loginid)
	if err != nil {
		log.Println(" User表更新错误 ： ",err)
		err = &e.Myerror{
			Code : e.ERROR_UPDATE_TABLE_CODE,
			Message : e.ERROR_UPDATE_TABLE_MESSAGE,
		}
		return
	}
	return
}
