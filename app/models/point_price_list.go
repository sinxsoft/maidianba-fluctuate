package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type PointPriceList struct {
	Id              int64 `orm:"pk"`
	PriceDate       string
	PriceDateAssist string
	Price           float64
	LastDayPrice    float64
	Mark            string
	HandleStatus    int
	/**
	来源  manual和auto
	注：手动插入的还是自动计算的，自动计算的是
	计算好插入的，设置状态handlestatus=1为已经处理
	*/
	Source string
}

func (p *PointPriceList) TableName() string {
	return TableName("point_price_list")
}

func (p *PointPriceList) Update(ormer orm.Ormer, fields ...string) error {
	if _, err := ormer.Update(p, fields...); err != nil {
		return err
	}
	return nil
}

func (p *PointPriceList) Add(ormer orm.Ormer) (int64, error) {
	return ormer.Insert(p)
}

func GetPointPriceListByTag(priceDate string, priceDateAssist string) (*PointPriceList, error) {
	c := new(PointPriceList)
	var pointPriceLists []PointPriceList
	result, errAll := orm.NewOrm().QueryTable(TableName("point_price_list")).
		Filter("priceDate", priceDate).Filter("priceDateAssist", priceDateAssist).All(&pointPriceLists)
	if errAll != nil {
		beego.Error(errAll)
		return nil, errAll
	}
	beego.Info("数量：" + strconv.FormatInt(result, 2))

	if len(pointPriceLists) > 0 {
		c = &pointPriceLists[0]
	} else {
		c = nil
	}

	//err := orm.NewOrm().QueryTable(TableName("point_price_list")).
	//	Filter("priceDate", priceDate).Filter("priceDateAssist", priceDateAssist).One(c)
	//if err != nil {
	//	return nil, err
	//}
	return c, nil
}
