package monitor

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func Check_home() (bool, error) {
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
	query.Add("station_id'", Conf.StationId)
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
	defer resp.Body.Close()
	if strings.Contains(string(body), BOOKABLE) {

		return true, nil
	} else {
		return false, errors.New("当前运力紧张")
	}

}
