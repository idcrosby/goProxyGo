package goProxy

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"github.com/idcrosby/web-tools"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var InfoLog *log.Logger

func init() {
	// init loggers
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
}

type GoProxy struct {
	// todo implement...
}

type RequestModifier interface {
	Modify(http.Request)
}

var DefaultGoProxy = &GoProxy{}


func (p *GoProxy) GoGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	// check(err)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	// check(err)
	return body, err
}

func (p *GoProxy) BuildRequest(url *url.URL, method string, body io.Reader, headers http.Header) *http.Request {
	req, err := http.NewRequest(method, url.String(), body)
	check(err)
	req.Header = headers
	return req
}

func (p *GoProxy) PreviewRequest(request *http.Request) {
	// request.Write()
}

func (p *GoProxy) ExecuteRequest(request *http.Request) (resp *http.Response, err error) {
	InfoLog.Println("Sending Request " + request.URL.String())
	return http.DefaultClient.Do(request)
}

func (p *GoProxy) Assault(request *http.Request, threads int, duration int) bool {
	start := time.Now().Unix()
	now := start;
	errors := make(chan error)
	// TODO Log data
	for now < (start + int64(duration)) {
		for i := 0; i < threads; i++ {
			// go ExecuteRequest(request)
			go func() {
				_, err := p.ExecuteRequest(request)
				errors <- err
			}()
		}
		now = time.Now().Unix()
	}
	// Check if any errors in channel
	for el := range errors {
		if el != nil {
			InfoLog.Println(el.Error())
			return false
		}
	}
	return true
}

func (p *GoProxy) GoGetAndFilter(url string, filter []string, pretty bool) []byte {

	body, err := p.GoGet(url)
	check(err)
	buf, err1 := myTools.JsonPositiveFilter(body, filter, pretty)
	check(err1)

	return buf
}

func (p *GoProxy) HateoasExpand(body []byte, expand []string) (buf []byte, err error) {
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
				expandField(node, sub)
			} else {
				node = node[sub].(map[string]interface{})
			}
		}
		// expandField(m, element)
		//delete(m, element)
	}
	buf, err = json.MarshalIndent(&m, "", "   ")
	check(err)
	return
}

func expandField(baseMap map[string]interface{}, el string) {
	fmt.Printf("expand field %s \n", el)
	// TODO handle potential casting error
	// var url string
	if subMap := baseMap[el].(map[string]interface{}); subMap != nil {
		// links := subMap["links"]
		// switch links.(type) {
		// 	case string:
		// 		url = links.(string)
		// 	case map[string]interface{}:
		// 		url = "object" //links.(map[string]interface{})["href"]
		// 	case []interface{}:
		// 		// subLink := links.([]interface{})[0]
		// 		url = "array" // subLink.(map[string]interface{})["href"]
		// 	default:
		// 		// TODO throw error
		// }

		// for initial implementation, only handle basic `href` field
		href := subMap["href"]
		// TODO this is base 64 encoding the response for some reason??
		baseMap[el], _ = DefaultGoProxy.GoGet(href.(string))
	}
}

func check(err error) { if err != nil { panic(err) } }