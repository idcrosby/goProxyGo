package goProxy

import (
	"io/ioutil"
	"github.com/idcrosby/web-tools"
	"net/http"
)

func GoGet(url string) string {
	res, err := http.Get(url)
	check(err)

	body, err := ioutil.ReadAll(res.Body)
	check(err)
	return string(body)
}

func GoGetAndFilter(url string, filter []string, pretty bool) []byte {
	res, err := http.Get(url)
	check(err)

	body, err := ioutil.ReadAll(res.Body)
	check(err)

	buf, err := myTools.JsonPositiveFilter(body, filter, pretty)
	check(err)

	return buf
}

func check(err error) { if err != nil { panic(err) } }