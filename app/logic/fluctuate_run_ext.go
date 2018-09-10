package logic

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"maidianba-fluctuate/app/models"
	"strconv"
	"time"
)

/**
自动计算，能进入这个方法说明PointPriceList今天没有记录
*/
func DoAutoFluctuate() error {
	currDate := time.Now().Format("2006-01-02")
	ppl, err := models.GetPointPriceListByTag(currDate, pointPriceAssistValue)
	if err != nil {
		beego.Error(err)
		return err
	}
	if ppl != nil {
		beego.Error(currDate + "已经存在记录，无法处理！")
		return nil
	}
	//通过公式生成记录插入
	lastDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	lastPointPriceList, errLast := models.GetPointPriceListByTag(lastDate, pointPriceAssistValue)
	if errLast != nil {
		beego.Error(err)
		return err
	}
	todayPrice, errFor := computeFormula(lastPointPriceList.Price)
	if errFor != nil {
		beego.Error(errFor)
		return errFor
	}
	beego.Info("今日价格：" + strconv.FormatFloat(todayPrice, 'E', -1, 64))
	newPointPriceList := new(models.PointPriceList)
	newPointPriceList.Id = time.Now().Unix()
	newPointPriceList.Price = todayPrice
	//没有处理，交给后续取处理
	newPointPriceList.HandleStatus = handleStatusNotOk
	newPointPriceList.LastDayPrice = lastPointPriceList.Price
	newPointPriceList.PriceDate = currDate
	newPointPriceList.PriceDateAssist = pointPriceAssistValue
	newPointPriceList.Source = autoComputeSource
	_, errorAdd := newPointPriceList.Add(orm.NewOrm())

	if errorAdd != nil {
		beego.Error(errorAdd)
		return errorAdd
	}
	//新开一个协程处理
	go func() {
		//交由手工计算程序去处理了
		errorManual := DoManualFluctuate(newPointPriceList)
		if errorManual != nil {
			beego.Error(errorManual)
		}
	}()

	return nil
}

func computeFormula(lastPrice float64) (float64, error) {
	pointFluctuateConfig, err := models.GetViewPointFluctuateConfigByTag()
	if err != nil {
		return 0, err
	}
	var price float64
	var a float64
	var b float64
	k2 := pointFluctuateConfig.Key2
	k1 := pointFluctuateConfig.Key1
	a = getValueForA()
	b = getValueForB()
	price = lastPrice * (1/k2 + (b/a)*(0.8/k2-0.86*k1-0.04))
	return price, nil
}

func getValueForA() float64 {
	return 1005
}

func getValueForB() float64 {
	return 1000
}
