package main

import (
	"dingdong_monitor/monitor"
	"fmt"
	"time"
)

func main() {

	monitor.InnitConfig("./")
	for true {
		res, err := monitor.Check_home()
		if res {
			monitor.PushToBark(monitor.NOTICE_TITLE, monitor.NOTICE_BOOKABLE, monitor.Conf.Bark.Sound)
		} else {
			fmt.Println(err)
		}
		if monitor.Conf.Mode == 1 {
			break
		}
		time.Sleep(time.Minute * 5)
	}

}
