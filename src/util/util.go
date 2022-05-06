package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Min(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}
func DownFile(remoteUrl string, savePath string) {

	// Get the data

	resp, err := http.Get(remoteUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		fmt.Println(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(out)

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}

}
