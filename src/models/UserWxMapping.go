package models

import (
	db "VPN-Manage/src/database"
	e "VPN-Manage/src/util/myerror"
	"fmt"
)

type UserWxMapping struct {
	Wxuuid	string  `xorm:"not null pk VARCHAR(32)"`
	Loginid  string `xorm:"not null VARCHAR(20)"`
}

func (u *UserWxMapping)Add() (err error) {
	_,err = db.Engine.Insert(u)
	if err != nil {
		errmsg := fmt.Sprintf("%s", err)
		err = &e.Myerror{
			Code : e.ERROR_INSERT_TABLE_CODE,
			Message : "表插入错误: " + errmsg,
		}
		return
	}
	return
}

func (u *UserWxMapping)Get() (mp UserWxMapping,err error) {
	ret,err := db.Engine.Alias("user_wx_mapping").Where("user_wx_mapping.wxuuid = ?", u.Wxuuid).Get(&mp)
	if err != nil{
		errmsg := fmt.Sprintf("%s", err)
		err = &e.Myerror{
			Code : e.ERROR_SELECT_TABLE_CODE,
			Message : "表查询错误" + errmsg,
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
/*
func (u *UserWxMapping)Update()(err error) {
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
} */
