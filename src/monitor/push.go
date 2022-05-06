package monitor

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func PushTo(title string, content string, sound string) {
	doPushToBark(title, content, sound)
}
func doPushToBark(title string, content string, sound string) {
	var urls []string
	for _, id := range Conf.Bark.Id {
		if id == "" || strings.ReplaceAll(id, " ", "") == "" {
			continue
		}
		u := "https://api.day.app/" + id + "/" + url.PathEscape(title) + "/" + url.PathEscape(content)
		urls = append(urls, u)
	}

	for _, u := range urls {
		req, err := http.NewRequest(http.MethodGet, u, nil)
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
