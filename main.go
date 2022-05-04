package main

import (
	monitor "dingdong_monitor/src/monitor"
	"fmt"
	"os"
	"time"
)

func main() {

	defer func() {
		fmt.Println("出现错误")
		time.Sleep(time.Second * 30)
		if monitor.Conf.Mode == 1 {
			os.Exit(-1)
		}
	}()
	monitor.InnitConfig("./")

	for true {
		now := time.Now()
		local := time.FixedZone("UTC+8", 8*60*60)

		fmt.Println("当前系统时间：" + now.In(local).Format("2006-01-02 15:04:05"))
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" "+"06:15:00", local)
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" "+"23:00:00", local)
		if !(now.In(local).Before(endTime) && now.In(local).After(startTime)) {
			fmt.Println("当前时间不在运行时间段内")
			if monitor.Conf.Mode == 1 {
				break
			} else {
				time.Sleep(time.Second * time.Duration(monitor.Conf.Rate))
				continue
			}
		}
		bookable, err := monitor.CheckHome()
		if bookable {
			fmt.Println("已可预约")
			monitor.PushToBark(monitor.NOTICE_TITLE, monitor.NOTICE_BOOKABLE, monitor.Conf.Bark.Sound)

		} else {
			fmt.Println(err)
		}
		if monitor.Conf.Mode == 1 {
			break
		}
		time.Sleep(time.Second * time.Duration(monitor.Conf.Rate))
	}

}
