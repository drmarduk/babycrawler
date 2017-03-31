package main

import (
	"io/ioutil"
	"net/http"
)

func downloadPage(url string) (string, error) {
	//log.Println("[+] download: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	src, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(src), nil
}
