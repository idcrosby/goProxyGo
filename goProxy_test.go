package goProxy

import(
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	url := "http://example.com"
	data := "..."
	if x, err := GoGet(url); x == nil || err != nil {
		t.Errorf("Get(" + url + ") = " + string(x) + ", want " + data)
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

func TestHateoasExpand(t *testing.T) {
	input := []byte("{\"one\":true,\"two\":{\"href\":\"http://api.icndb.com/jokes/15\"}}")
	// url := "http://api.icndb.com/jokes/15"
	expand := []string{"two"}
	data := "{\"one\":true,\"two\":{\"value\":{\"joke\":\"When Chuck Norris goes to donate blood, he declines the syringe, and instead requests a hand gun and a bucket.\"}}}"
	if x, err := HateoasExpand(input, expand); string(x) != data  || err != nil {
		t.Errorf("HateoasExpand(" + string(input) + ", " + arrayToString(expand) + ") = " + string(x) + ", want " + data)
	}	
}

func TestAssault(t *testing.T) {
	req, _ := http.NewRequest("GET",  "http://api.icndb.com/jokes/15", nil)
	if !Assault(req, 2, 3) {
		t.Errorf("Assault(http://api.icndb.com/jokes/15, 2, 3) = false")
	}
}


// Util Methods

func arrayToString(input []string) (output string) {
	for _, value := range input { output += string(value) }
	return
}