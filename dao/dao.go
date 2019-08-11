package dao

import (
	"proxy_manage/pojo"
	"proxy_manage/config"
)



//保存操作
func Save(proxy pojo.ProxyPool)  {
	db := config.GetDb()
    db.Save(&proxy)
}
//查询所有
func FindAll() []pojo.ProxyPool{
	db:=config.GetDb()
	var po []pojo.ProxyPool
	db.Find(&po)
	return po
}

func DeleteById(id string)  {
	db := config.GetDb()
	db.Model(&pojo.ProxyPool{}).Delete(pojo.ProxyPool{Id:id})
}

func SvaeDomain(daomain pojo.Domain)  {
	db := config.GetDb()
	db.Create(&daomain)
}
