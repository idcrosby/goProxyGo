package goProxy

import (
	"encoding/json"
	"io/ioutil"
	"github.com/idcrosby/web-tools"
	"net/http"
	"strings"
)

func GoGet(url string) []byte {
	res, err := http.Get(url)
	check(err)

	body, err := ioutil.ReadAll(res.Body)
	check(err)
	return body
}

func GoGetAndFilter(url string, filter []string, pretty bool) []byte {

	body := GoGet(url)
	buf, err := myTools.JsonPositiveFilter(body, filter, pretty)
	check(err)

	return buf
}

func HateoasExpand(body []byte, expand []string) (buf []byte, err error) {
 	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}
	// Access the data's underlying interface
	m := f.(map[string]interface{})

	for _,element := range expand {
		subEls := strings.Split(element,".")
		node := m
		for index,sub := range subEls {
			//Check if last element
			if (index >= (len(subEls) -1)) {
				//delete(node, sub)
				expandThis(node, sub)
			} else {
				node = node[sub].(map[string]interface{})
			}
		}
		expandThis(m, element)
		//delete(m, element)
	}
	buf, err = json.MarshalIndent(&m, "", "   ")
	check(err)
	return
}

func expandThis(baseMap map[string]interface{}, el string) {
	// TODO handle potential casting error
	var url string
	if subMap := baseMap[el].(map[string]interface{}); subMap != nil {
		links := subMap["links"]
		switch links.(type) {
			case string:
				url = links.(string)
			case map[string]interface{}:
				url = "object" //links.(map[string]interface{})["href"]
			case []interface{}:
				// subLink := links.([]interface{})[0]
				url = "array" // subLink.(map[string]interface{})["href"]
			default:
				// TODO throw error
		}
		baseMap[el] = GoGet(url)
	}
}

func check(err error) { if err != nil { panic(err) } }