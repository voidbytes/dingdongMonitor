package monitor

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func PushToBark(title string, content string, sound string) {
	var urls []string
	for _, id := range Conf.Bark.Id {
		url := "https://api.day.app/" + id + "/" + title + "/" + content
		urls = append(urls, url)
	}

	for _, url := range urls {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}
		if resp.StatusCode != 200 {
			fmt.Println("请检查bark是否配置正确")
		}
	}

}
