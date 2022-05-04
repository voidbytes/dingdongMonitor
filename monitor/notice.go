package monitor

import (
	"crypto/tls"
	"errors"
	"net/http"
)

func PushToBark(title string, content string, sound string) error {
	url := "https://api.day.app/" + Conf.Bark.Id + "/" + title + "/" + content
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	var client = &http.Client{
		Timeout:   TIME_OUT,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	query := req.URL.Query()
	query.Add("isArchive", "1")
	if sound != "" {
		query.Add("sound", sound)

	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("请检查bark是否配置正确")
	}
	return nil

}
