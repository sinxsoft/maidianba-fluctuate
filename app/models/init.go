package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/url"
)

type BitType int

var (
	BitBoolFalse BitType = 0
	BitBoolTrue  BitType = 1
)

//为什么不用  _ 引用 取调用init方法呢？因为这个Init保证了调用顺序
func Init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)

	orm.RegisterModel(new(PointFluctuate))
	orm.RegisterModel(new(PointPriceList))
	orm.RegisterModel(new(Task))
	orm.RegisterModel(new(TaskLog))
	orm.RegisterModel(new(PointFluctuateConfig))

	orm.RegisterModel(new(ViewPointFluctuateConfig))
	orm.RegisterModel(new(ViewPointFluctuate))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}

func CreateOrmer() orm.Ormer {
	return orm.NewOrm()
}
