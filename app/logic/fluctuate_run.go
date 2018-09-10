package logic

import (
	"fmt"
	"github.com/astaxie/beego"
	"maidianba-fluctuate/app/models"
	"time"
)

const pointPriceAssistValue = "00:00:00"
const handleStatusOk = 1
const handleStatusNotOk = 0
const autoComputeTrue = 1
const autoComputeFalse = 0
const autoComputeSource = "auto"

type FluctuateHandler struct {
	Handler
}

/*
执行计算买点波动的逻辑：
1，手动的话，按照point_price_list 波动
2，自动的话，按照公式波动
根据配置来决定到底什么时候开始按照2执行
*/
func (fh *FluctuateHandler) Handle() error {

	beego.Info("执行开始...")

	autoCompute := autoComputeFalse
	pointFluctuateConfig, err := models.GetViewPointFluctuateConfigByTag()
	if err == nil {
		autoCompute = pointFluctuateConfig.IsAutoCompute
	}

	if autoCompute == autoComputeTrue {
		//需要通过公式（配置到数据库中autoCompute=1）计算,
		//那么就通过公式来运行
		return DoAutoFluctuate()
	} else {
		currDate := time.Now().Format("2006-01-02")
		ppl, err := models.GetPointPriceListByTag(currDate, pointPriceAssistValue)
		if err != nil {
			beego.Error(err)
			return err
		}

		if ppl.HandleStatus == handleStatusOk {
			beego.Error(currDate + "已经处理过了！")
			return nil
		}
		return DoManualFluctuate(ppl)
	}

	beego.Info("执行结束！")
	return nil
}

/**
手动计算
*/
func DoManualFluctuate(pointPriceList *models.PointPriceList) error {
	//查询list，时间为当天的记录
	//如果存在，就覆盖，如果不存在
	//就沿用头一天的价格
	pointFluctuate, err := models.GetViewPointFluctuateByIsUse(true)
	if err != nil {
		beego.Error(err)
		return err
	}
	if pointFluctuate == nil {
		beego.Error("pointFluctuate对象isUse为1的对象不存在")
		return nil
	}

	var endTime time.Time

	endTime = time.Now()
	//创建统一的ormer，保证一个事务
	ormer := models.CreateOrmer()
	//开始事务
	ormer.Begin()

	pointFluctuate.UpdateTime = time.Now()
	pointFluctuate.IsUse = 0
	pointFluctuate.EndDate = endTime
	errUpdate := pointFluctuate.Update(ormer, "isUse", "endDate", "updateTime")

	if errUpdate != nil {
		beego.Error("失败： " + errUpdate.Error())
		ormer.Rollback()
		return nil
	}

	//新插入的对象
	pointFluctuateNew := new(models.PointFluctuate)
	pointFluctuateNew.Id = time.Now().Unix()
	//pointFluctuateNew.Id = time.Time{}.Unix()
	pointFluctuateNew.BeginDate = endTime
	pointFluctuateNew.IsUse = 1
	pointFluctuateNew.Remark = "手动生成"
	pointFluctuateNew.CreateTime = time.Now()
	pointFluctuateNew.CreateUser = "直接任务"
	pointFluctuateNew.IsDelete = 0

	//pointFluctuateNew.CurrentPrice = "2.99"
	//pointFluctuateNew.Balance = "0.09"
	//pointFluctuateNew.PriceRate = "2.3"
	//调用三个值得设置方法
	errSet := setPointFluctuateNewValue(pointFluctuateNew, pointPriceList)
	if errSet != nil {
		pointFluctuateNew.CurrentPrice = pointFluctuate.CurrentPrice
		pointFluctuateNew.Balance = pointFluctuate.Balance
		pointFluctuateNew.PriceRate = pointFluctuate.PriceRate
		pointFluctuateNew.Remark = "生成出错"
	}

	_, errAdd := pointFluctuateNew.Add(ormer)
	if errAdd == nil {
		beego.Info("新价格生产成功！")
	} else {
		beego.Error("新价格失败！详情：" + errAdd.Error())
		ormer.Rollback()
		return errAdd
	}

	//处理point_price_list表，置HandleStatus为已经处理
	pointPriceList.HandleStatus = handleStatusOk
	errPointPriceList := pointPriceList.Update(ormer, "handleStatus")

	if errPointPriceList == nil {
		beego.Info("handleStatus修改为1成功！")
	} else {
		beego.Error("handleStatus修改为1失败！详情：" + errPointPriceList.Error())
		ormer.Rollback()
		return errPointPriceList
	}

	ormer.Commit()
	beego.Info(endTime.String() + "事务已经提交！")
	return nil
}

func setPointFluctuateNewValue(pf *models.PointFluctuate, pointPriceList *models.PointPriceList) error {
	//currDate := time.Now().Format("2006/01/02")
	//ppl, err := models.GetPointPriceListByTag(currDate, "00:00:00")
	//if err != nil {
	//ppl.LastDayPrice

	var currPrice float64
	var balance float64
	var priceRate float64

	currPrice = pointPriceList.Price
	balance = pointPriceList.Price - pointPriceList.LastDayPrice
	priceRate = ((pointPriceList.Price - pointPriceList.LastDayPrice) / pointPriceList.LastDayPrice) * 100

	//fmt.Sprintf("%0.6f", 17.82671567890123456789987654324567898765432)

	pf.CurrentPrice = fmt.Sprintf("%0.6f", currPrice)
	pf.Balance = fmt.Sprintf("%0.2f", balance)
	pf.PriceRate = fmt.Sprintf("%0.2f", priceRate)
	return nil
	//} else {
	//	return err
	//}
}
