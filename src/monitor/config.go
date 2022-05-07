package monitor

import (
	"dingdong_monitor/src/util"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	UA              = "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/19E258 Ariver/1.1.0 AliApp(AP/10.2.60.6200) Nebula WK RVKType(0) AlipayDefined(nt:WIFI,ws:390|780|3.0,ac:T) AlipayClient/10.2.60.6200 Language/zh-Hans Region/CN NebulaX/1.0.0"
	CITY            = "0101"
	API_VERSION     = "9.51.0"
	BUILD_VERSION   = "2.86.3"
	TIME_OUT        = 10 * time.Second
	BOOKABLE        = "可预约"
	NOTICE_BOOKABLE = "可以预约啦"
	NOTICE_TITLE    = "运力监控"
	DEFAULT_RATE    = 3600
	LOCAL_MODE      = 0
	GITHUB_MODE     = 1
	APP_CLIENT_ID   = 10
	ALIMINIMARK     = "DDNMSL"
	OPEN_ID         = "6666666666666666"
	REFERER         = "https://2021001157662937.hybrid.alipay-eco.com/2021001157662937/0.2.2205051588.15/index.html#pages/mainPackage/home/home"
)

var Conf = new(config)

type Keyword struct {
	Name  string `mapstructure:"name"`
	Price string `mapstructure:"price"`
}
type config struct {
	Mode                 int    `mapstructure:"mode"`
	Rate                 uint   `mapstructure:"rate"`
	IsEnableStockMonitor bool   `mapstructure:"stock_monitor"`
	IsPrivate            bool   `mapstructure:"private"`
	StationId            string `mapstructure:"station_id"`
	Longitude            string `mapstructure:"longitude"`
	Latitude             string `mapstructure:"latitude"`

	KeyWords []Keyword `mapstructure:"keywords"`

	Bark struct {
		Id    []string `mapstructure:"id"`
		Sound string   `mapstructure:"sound"`
	}
}

func InnitConfig() {
	var conf = new(config)
	argConf, configPath := ConfigFromArg()
	fileConf := ConfigFromFile(configPath)
	conf = &fileConf
	if argConf.Latitude != "" && argConf.Longitude != "" {
		conf.Longitude = argConf.Longitude
		conf.Latitude = argConf.Latitude
	}
	if argConf.StationId != "" {
		conf.StationId = argConf.StationId
	}
	if argConf.Mode != -1 {
		conf.Mode = argConf.Mode
	}
	if len(argConf.Bark.Id) != 0 {
		set := make(map[string]byte)
		for _, id := range argConf.Bark.Id {
			if id != "" {
				set[id] = 0
			}
		}
		for _, id := range conf.Bark.Id {
			if id != "" {
				set[id] = 0
			}
		}
		var temp []string
		for k, _ := range set {
			temp = append(temp, k)

		}
		conf.Bark.Id = temp
	}
	if conf.StationId == "" {
		fmt.Println("没有填写 station_id，尝试从经纬度获取")
		if conf.Latitude == "" || conf.Longitude == "" {
			panic(fmt.Errorf("请至少填写station_id和经纬度的其中一项"))
		}
		conf.StationId = GetStationId(conf.Longitude, conf.Latitude)
		if conf.StationId == "" {
			panic(fmt.Errorf("获取站点信息失败"))
		}

	}
	if len(conf.Bark.Id) == 0 {
		fmt.Println("require bark_id")
		panic(fmt.Errorf("require bark_id"))
	}
	if conf.Mode == LOCAL_MODE {
		fmt.Println("当前为本机模式运行")
	} else if conf.Mode == 1 {
		fmt.Println("当前为GitHub Action模式运行")
		fmt.Println("GitHub模式下config.yaml文件里的频率设置均不生效  需要在GitHub action配置文件配置")
	}
	if conf.Rate == 0 {
		conf.Rate = DEFAULT_RATE
	}
	Conf = conf

}

type barkIds []string

//实现String接口
func (i *barkIds) String() string {
	return fmt.Sprintf("%s", *i)
}
func (i *barkIds) Set(value string) error {
	//反正后面会去重的
	if len(*i) > 0 {
		for _, id := range strings.Split(value, ",") {
			barkId := id

			*i = append(*i, barkId)
		}
	}
	for _, id := range strings.Split(value, ",") {
		barkId := id

		*i = append(*i, barkId)
	}
	return nil
}

func ConfigFromArg() (*config, string) {
	var conf = new(config)
	configFilePath := flag.String("p", "./config.yaml", "配置文件路径指定")
	mode := flag.Int("m", -1, "mode")
	rate := flag.Uint("r", 3600, "rate")
	lat := flag.String("lat", "", "latitude")
	lng := flag.String("lng", "", "lng")
	isPrivate := flag.Bool("i", false, "private info")
	stationId := flag.String("sta", "", "station id")
	//stockMonitorFlag:=flag.Bool("s",false,"stock monitor")
	var barkIdsFlag barkIds
	flag.Var(&barkIdsFlag, "b", "bark id")
	flag.Parse()
	conf.Mode = *mode
	conf.Rate = *rate
	conf.Latitude = *lat
	conf.Longitude = *lng
	conf.StationId = *stationId
	conf.Bark.Id = barkIdsFlag
	if *isPrivate {
		fmt.Println("隐私模式")
	}
	conf.IsPrivate = *isPrivate
	Conf.IsPrivate = *isPrivate
	return conf, *configFilePath

}
func ConfigFromFile(path string) config {
	var conf = new(config)
	if !Conf.IsPrivate {
		fmt.Println("配置文件将会从 " + path + " 加载")
	}

	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		fmt.Println("配置文件下载中")
		util.DownFile(path, "./config.remote.yaml")
		path = "./config.remote.yaml"
		fmt.Println("配置文件下载完毕")

	}

	viper.SetConfigFile(path)
	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("fatal error reading configuration file: %s \n", err)
		panic(fmt.Errorf("fatal error reading configuration file: %s \n", err))
	}
	// Unmarshal configuration
	if err := viper.Unmarshal(conf); err != nil {
		fmt.Printf("unmarshal configuration failed, err: %s \n", err)
		panic(fmt.Errorf("unmarshal configuration failed, err: %s \n", err))
	}
	return *conf

}
