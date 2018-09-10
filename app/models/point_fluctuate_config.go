package models

import "github.com/astaxie/beego/orm"

//`id` int(11) NOT NULL,
//`key1` decimal(12,6) NOT NULL,
//`key2` decimal(12,6) NOT NULL,
//`key_ext1` decimal(12,6) DEFAULT NULL,
//`is_validated` bit(1) NOT NULL,
//`is_auto_compute` bit(1) NOT NULL
type PointFluctuateConfig struct {
	Id            int64
	Key1          float64
	Key2          float64
	KeyExt1       float64
	IsValidated   int
	IsAutoCompute int
}

func (p *PointFluctuateConfig) TableName() string {
	return TableName("point_fluctuate_config")
}

type ViewPointFluctuateConfig struct {
	//Id int64
	PointFluctuateConfig
}

func (p *ViewPointFluctuateConfig) TableName() string {
	return TableName("view_point_fluctuate_config")
}

func GetViewPointFluctuateConfigByTag() (*ViewPointFluctuateConfig, error) {
	c := new(ViewPointFluctuateConfig)
	err := orm.NewOrm().QueryTable(TableName("view_point_fluctuate_config")).
		Filter("id", 1).One(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetPointFluctuateConfigByTag() (*PointFluctuateConfig, error) {
	c := new(PointFluctuateConfig)
	err := orm.NewOrm().QueryTable(TableName("view_point_fluctuate_config")).
		Filter("id", 1).One(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
