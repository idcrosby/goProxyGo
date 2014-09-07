package goProxy

import(
	"testing"
)

func TestGet(t *testing.T) {
	url := "http://example.com"
	data := "..."
	if x := GoGet(url); len(x) == 0 {
		t.Errorf("Get(" + url + ") = " + x + ", want " + data)
	}	
}

func TestFilterGet(t *testing.T) {
	url := "http://api.icndb.com/jokes/15"
	filter := []string{"value.joke"}
	data := "{\"value\":{\"joke\":\"When Chuck Norris goes to donate blood, he declines the syringe, and instead requests a hand gun and a bucket.\"}}"
	if x := GoGetAndFilter(url, filter, false); string(x) != data  {
		t.Errorf("GoGetAndFilter(" + url + ") = " + string(x) + ", want " + data)
	}	
}