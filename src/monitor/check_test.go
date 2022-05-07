package monitor

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"testing"
)

func TestGetStationId(t *testing.T) {
	_ = GetStationId("121.547176", "31.236673")

}
func TestCheckStock(t *testing.T) {
	conf := ConfigFromFile("../../config.yaml")
	Conf = &conf
	Conf.StationId = GetStationId("121.547176", "31.236673")

	keywords := Conf.KeyWords
	var newKeywords []Keyword
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
	ok, _, _, total := CheckStock(1, newKeywords)
	if !ok {
		panic("error")
	}
	fmt.Println("total:" + strconv.Itoa(total))

}
func TestCheckTransportCapacity(t *testing.T) {
	conf := ConfigFromFile("../../config.yaml")
	Conf = &conf
	Conf.StationId = GetStationId("121.547176", "31.236673")
	_, err := CheckTransportCapacity()
	if err != nil {
		fmt.Println(err)
	}
}
