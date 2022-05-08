package monitor

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func CheckTransportCapacity() (bool, error) {
	url := "https://maicai.api.ddxq.mobi/homeApi/newDetails"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {

		return false, err
	}
	var client = &http.Client{
		Timeout:   TIME_OUT,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	req.Header.Add("user-agent", UA)
	req.Header.Add("alipayminimark", ALIMINIMARK)
	req.Header.Add("referer", REFERER)
	req.Header.Add("ddmc-app-client-id", strconv.Itoa(APP_CLIENT_ID))
	req.Header.Add("ddmc-city-number", CITY)
	req.Header.Add("ddmc-api-version", API_VERSION)
	req.Header.Add("ddmc-build-version", BUILD_VERSION)

	query := req.URL.Query()
	query.Add("openid", OPEN_ID)
	query.Add("api_version", API_VERSION)
	query.Add("app_version", BUILD_VERSION)
	query.Add("station_id", Conf.StationId)
	query.Add("city_number", CITY)
	query.Add("buildVersion", BUILD_VERSION)
	query.Add("app_client_id", strconv.Itoa(APP_CLIENT_ID))
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 || !strings.Contains(string(body), "成功") {
		return false, errors.New(string(body))
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	if strings.Contains(string(body), BOOKABLE) {
		fmt.Println("当前有运力")
		return true, nil
	} else {
		return false, errors.New("当前运力紧张，或者首页暂无信息")
	}

}
func GetStationId(lng string, lat string) string {
	fmt.Println("正在获取站点信息...")
	url := "https://sunquan.api.ddxq.mobi/api/v2/user/location/refresh/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var client = &http.Client{
		Timeout:   TIME_OUT,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	req.Header.Add("user-agent", UA)
	req.Header.Add("ddmc-longitude", lng)
	req.Header.Add("ddmc-latitude", lat)
	req.Header.Add("ddmc-app-client-id", strconv.Itoa(APP_CLIENT_ID))
	query := req.URL.Query()
	query.Add("api_version", API_VERSION)
	query.Add("station_id", Conf.StationId)
	query.Add("city_number", CITY)
	query.Add("buildVersion", BUILD_VERSION)
	query.Add("app_client_id", strconv.Itoa(APP_CLIENT_ID))
	query.Add("longitude", lng)
	query.Add("latitude", lat)
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	data := result["data"].(map[string]interface{})
	StationId := data["station_id"].(string)
	if StationId != "" {
		info := data["station_info"].(map[string]interface{})
		name := info["name"].(string)
		if !Conf.IsPrivate {
			fmt.Println("当前经纬度读取到的站点名称: " + name)
		}

	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	return StationId

}
func CheckStock(page int, keyWords []Keyword) (isSucccess bool, isMore bool, products []map[string]string, totalGoods int) {
	url := "https://maicai.api.ddxq.mobi/homeApi/homeFlowDetail"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return false, false, nil, 0
	}
	var client = &http.Client{
		Timeout:   TIME_OUT,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	req.Header.Add("user-agent", UA)
	req.Header.Add("alipayminimark", ALIMINIMARK)
	req.Header.Add("referer", REFERER)
	req.Header.Add("ddmc-app-client-id", strconv.Itoa(APP_CLIENT_ID))
	req.Header.Add("ddmc-city-number", CITY)
	req.Header.Add("ddmc-api-version", API_VERSION)
	req.Header.Add("ddmc-build-version", BUILD_VERSION)
	req.Header.Add("ddmc-station-id", Conf.StationId)
	query := req.URL.Query()
	query.Add("api_version", API_VERSION)
	query.Add("app_version", BUILD_VERSION)
	query.Add("station_id", Conf.StationId)
	query.Add("city_number", CITY)
	query.Add("buildVersion", BUILD_VERSION)
	query.Add("app_client_id", strconv.Itoa(APP_CLIENT_ID))
	query.Add("s_id", OPEN_ID)
	query.Add("open_id", OPEN_ID)
	query.Add("longitude", Conf.Longitude)
	query.Add("latitude", Conf.Latitude)
	query.Add("page", strconv.Itoa(page))
	query.Add("tab_type", "1")
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, false, nil, 0
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false, false, nil, 0
	}

	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return false, false, nil, 0
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)
	data := result["data"].(map[string]interface{})
	isMore = data["is_more"].(bool)
	productsList := data["list"].([]interface{})

	products = make([]map[string]string, 0)
	totalGoods = 0
	for _, product := range productsList {
		status := product.(map[string]interface{})["status"]
		if status == nil {
			continue
		}
		if int(status.(float64)) != 1 {
			continue
		}
		totalGoods = totalGoods + 1
		for _, keyword := range keyWords {

			name := product.(map[string]interface{})["name"]
			price := product.(map[string]interface{})["price"]
			setPrice, err := decimal.NewFromString(keyword.Price)
			if err != nil {
				fmt.Println(err)
			}
			objPrice, err := decimal.NewFromString(price.(string))
			if err != nil {
				fmt.Println(err)
			}
			if !strings.Contains(name.(string), keyword.Name) {
				continue

			}
			if !setPrice.Sub(objPrice).IsPositive() && !setPrice.Equal(decimal.NewFromInt(0)) {
				continue

			}
			res := map[string]string{
				"keyword": keyword.Name,
				"name":    name.(string),
				"price":   price.(string),
			}
			products = append(products, res)

		}
	}
	if len(products) != 0 {
		return true, isMore, products, totalGoods

	}

	return true, isMore, nil, totalGoods

}
