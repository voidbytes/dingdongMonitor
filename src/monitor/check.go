package monitor

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func CheckHome() (bool, error) {
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
	query := req.URL.Query()
	query.Add("api_version", API_VERSION)
	query.Add("station_id", Conf.StationId)
	query.Add("city_number", CITY)
	query.Add("buildVersion", BUILD_VERSION)
	query.Add("app_client_id", "1")
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
		return true, nil
	} else {
		return false, errors.New("当前运力紧张，或者首页暂无信息")
	}

}
func GetStationId() string {
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
	req.Header.Add("ddmc-longitude", Conf.Longitude)
	req.Header.Add("ddmc-latitude", Conf.Latitude)
	query := req.URL.Query()
	query.Add("api_version", API_VERSION)
	query.Add("station_id", Conf.StationId)
	query.Add("city_number", CITY)
	query.Add("buildVersion", BUILD_VERSION)
	query.Add("app_client_id", "1")
	query.Add("longitude", Conf.Longitude)
	query.Add("latitude", Conf.Latitude)
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
		fmt.Println("当前经纬度读取到的站点名称: " + name)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	return StationId

}
