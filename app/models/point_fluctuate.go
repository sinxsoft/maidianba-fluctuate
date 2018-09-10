package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//`id` bigint(20) NOT NULL COMMENT '主键id',
//`current_price` varchar(50) DEFAULT NULL COMMENT '当前价格',
//`balance` varchar(50) DEFAULT NULL COMMENT '差额',
//`price_rate` varchar(50) DEFAULT NULL COMMENT '增长率',
//`begin_date` datetime DEFAULT NULL COMMENT '开始时间',
//`end_date` datetime DEFAULT NULL COMMENT '结束时间',
//`is_use` bit(1) DEFAULT b'1' COMMENT '是否使用，1-正在使用，0-停止使用',
//`remark` varchar(500) DEFAULT NULL COMMENT '备注',
//`create_user` varchar(100) DEFAULT NULL COMMENT '创建人',
//`create_time` datetime DEFAULT NULL COMMENT '创建时间',
//`update_user` varchar(100) DEFAULT NULL COMMENT '更新人',
//`update_time` datetime DEFAULT NULL COMMENT '更新时间',
//`ext1` varchar(100) DEFAULT NULL COMMENT '扩展1',
//`ext2` varchar(100) DEFAULT NULL COMMENT '扩展2',
//`ext3` varchar(100) DEFAULT NULL COMMENT '扩展3',
//`is_delete` int(2) DEFAULT '0' COMMENT '是否删除:0、否 1、是',
type PointFluctuate struct {
	Id           int64
	CurrentPrice string
	Balance      string
	PriceRate    string
	BeginDate    time.Time
	EndDate      time.Time
	IsUse        int
	Remark       string
	CreateUser   string
	CreateTime   time.Time
	UpdateUser   string
	UpdateTime   time.Time
	Ext1         string
	Ext2         string
	Ext3         string
	IsDelete     int
}

func (p *PointFluctuate) TableName() string {
	return TableName("point_fluctuate")
}

type ViewPointFluctuate struct {
	//Id int64
	PointFluctuate
}

func (p *ViewPointFluctuate) TableName() string {
	return TableName("view_point_fluctuate")
}

func (p *PointFluctuate) Update(ormer orm.Ormer, fields ...string) error {
	if _, err := ormer.Update(p, fields...); err != nil {
		return err
	}
	return nil
}

func (p *PointFluctuate) Add(ormer orm.Ormer) (int64, error) {
	return ormer.Insert(p)
}

func GetViewPointFluctuateByIsUse(isUse bool) (*ViewPointFluctuate, error) {
	flag := 0
	if isUse {
		flag = 1
	}
	pf := new(ViewPointFluctuate)
	err := orm.NewOrm().QueryTable(TableName("view_point_fluctuate")).Filter("isUse", flag).One(pf)
	if err == nil {
		return pf, nil
	} else {
		return nil, err
	}
}

func GetPointFluctuateByIsUse(isUse bool) (*PointFluctuate, error) {
	flag := 0
	if isUse {
		flag = 1
	}
	pf := new(PointFluctuate)
	err := orm.NewOrm().QueryTable(TableName("view_point_fluctuate")).Filter("isUse", flag).One(pf)
	if err == nil {
		return pf, nil
	} else {
		return nil, err
	}
}
