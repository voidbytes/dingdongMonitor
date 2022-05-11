package main

import (
	monitor "dingdong_monitor/src/monitor"
	"fmt"
	"github.com/shopspring/decimal"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {

	defer func() {
		if monitor.Conf.Mode == monitor.LOCAL_MODE {
			fmt.Println("出现错误")
			time.Sleep(time.Second * 30)
			if monitor.Conf.Mode == 1 {
				os.Exit(-1)
			}
		}
	}()
	monitor.InnitConfig()

	for true {

		if !isAllowTimeNow() {
			fmt.Println("当前时间不在运行时间段内")
			if monitor.Conf.Mode == monitor.GITHUB_MODE {
				//github模式直接退出 等待下一次GitHub action调用 避免运行时间过长被GitHub风控
				break
			} else {
				if monitor.Conf.IsEnableStockMonitor {
					time.Sleep(time.Duration(monitor.Conf.Rate) * time.Second)
				}
				continue
			}
		}

		if monitor.Conf.IsEnableStockMonitor {
			CheckStockAndTrans()
		} else {
			CheckTrans()
		}

		if monitor.Conf.Mode == monitor.GITHUB_MODE {
			break
		}
		time.Sleep(time.Duration(monitor.Conf.Rate) * time.Second)
	}

}
func CheckTrans() bool {
	bookable, err := monitor.CheckTransportCapacity()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if bookable && !monitor.Conf.IsEnableStockMonitor {
		fmt.Println("已可预约")
		monitor.PushTo(monitor.NOTICE_TITLE, monitor.NOTICE_BOOKABLE, monitor.Conf.Bark.Sound)
		return true
	}
	if bookable {
		return true
	}
	return false

}
func CheckStockAndTrans() bool {
	bookable := CheckTrans()
	if !bookable {
		return false
	}

	keywords := monitor.Conf.KeyWords
	var newKeywords []monitor.Keyword
	for _, keyword := range keywords {
		if keyword.Name == "" {
			continue

		}
		d, err := decimal.NewFromString(keyword.Price)
		if err != nil || d.IsNegative() {
			continue
		}

		newKeywords = append(newKeywords, keyword)

	}
	page := 1
	for true {
		isSuccess, isMore, products, totalGoods := monitor.CheckStock(page, newKeywords)
		page = page + 1

		if isSuccess {
			if len(newKeywords) == 0 && totalGoods > 0 {
				fmt.Println("首页有商品可以购买")
				monitor.PushTo("当前有物品可以购买", "详情请查看App", monitor.Conf.Bark.Sound)
				return true
			}
			for _, product := range products {
				reg := regexp.MustCompile("\\s+")
				keyword := product["keyword"]
				name := product["name"]
				name = reg.ReplaceAllString(name, "")
				price := product["price"]
				time.Sleep(100 * time.Microsecond)
				fmt.Println("您关注的 " + " 已可购买")
				fmt.Println("商品名: " + name)
				fmt.Println("价格:" + price)
				monitor.PushTo(strings.Trim("您关注的"+string([]rune(keyword)[:6]), "\x00")+"已可购买", "商品名:"+strings.Trim(string([]rune(name)[:10]), "\x00")+" 价格:"+price, monitor.Conf.Bark.Sound)

			}
			if !isMore {
				break
			}
		}
	}
	return true
}
func isAllowTimeNow() bool {
	now := time.Now()
	local := time.FixedZone("UTC+8", 8*60*60)

	fmt.Println("当前系统时间：" + now.In(local).Format("2006-01-02 15:04:05"))
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", now.In(local).Format("2006-01-02")+" "+"06:15:00", local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", now.In(local).Format("2006-01-02")+" "+"23:00:00", local)
	if now.In(local).Before(endTime) && now.In(local).After(startTime) {
		return true
	}
	return false
}
