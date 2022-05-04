package monitor

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	UA              = "neighborhood/9.50.2 (iPhone; iOS 15.4.1; Scale/3.00)"
	CITY            = "0101"
	API_VERSION     = "9.50.26"
	BUILD_VERSION   = "1232"
	TIME_OUT        = 10 * time.Second
	BOOKABLE        = "可预约"
	NOTICE_BOOKABLE = "可以预约啦"
	NOTICE_TITLE    = "运力监控"
)

var Conf = new(config)

type config struct {
	Mode      int    `mapstructure:"mode"`
	StationId string `mapstructure:"station_id"`
	Bark      struct {
		Id    string `mapstructure:"id"`
		Sound string `mapstructure:"sound"`
	}
}

func InnitConfig(path string) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading configuration file: %s \n", err))
	}
	// Unmarshal configuration
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal configuration failed, err: %s \n", err))
	}

	if Conf.StationId == "" {
		panic(fmt.Errorf("rquire station_id"))
	}
	if Conf.Bark.Id == "" {
		panic(fmt.Errorf("require bark_id"))
	}
}
