package form

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestForm(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com?user_id=180&name=tim", nil)
	if err != nil {
		t.Fatal(err)
	}

	var userID int
	var name string

	form := Form{
		Fields: []*Field{
			{Name: "user_id", Value: &IntField{&userID}},
			{Name: "name", Value: &StringField{&name}},
		},
	}

	parseErr := form.Parse(request)

	if parseErr != nil {
		t.Fatal(parseErr)
	}

	if 180 != userID {
		t.Errorf("did not load integer; want %d, got %d", 180, userID)
	}

	if "tim" != name {
		t.Errorf("did not load string; want %q, got %q", "tim", name)
	}
}

func TestFormIntParseError(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com?user_id=nowayjose", nil)
	if err != nil {
		t.Fatal(err)
	}

	var userID int

	form := Form{
		Fields: []*Field{
			{Name: "user_id", Value: &IntField{&userID}},
		},
	}

	parseErr := form.Parse(request)

	if parseErr == nil {
		t.Error("Expected to get an error when parsing the user_id")
	}

	if 0 != userID {
		t.Errorf("did not load integer; want %d, got %d", 180, userID)
	}
}

func TestRequiredParameters(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	var userID int

	form := Form{
		Fields: []*Field{
			{Name: "user_id", Value: &IntField{&userID}, Required: true},
		},
	}

	parseErr := form.Parse(request)
	if parseErr == nil {
		t.Errorf("expected a required params error")
	}
}

func TestRequiredStringParameters(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	var name string

	form := Form{
		Fields: []*Field{
			{Name: "name", Value: &StringField{&name}, Required: true},
		},
	}

	parseErr := form.Parse(request)
	if parseErr == nil {
		t.Errorf("expected a required params error")
	}
}

func TestOptionStyleInterface(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com?user_id=180&name=tim", nil)
	if err != nil {
		t.Fatal(err)
	}

	var userID int
	var name string

	form := NewForm()

	form.AddField("user_id", FieldValue(&IntField{&userID}), QueryParam, Required)
	form.AddField("name", FieldValue(&StringField{&name}), QueryParam, NotRequired)

	parseErr := form.Parse(request)

	if parseErr != nil {
		t.Fatal(parseErr)
	}

	if 180 != userID {
		t.Errorf("did not load integer; want %d, got %d", 180, userID)
	}

	if "tim" != name {
		t.Errorf("did not load string; want %q, got %q", "tim", name)
	}

}

func TestFormValueParams(t *testing.T) {
	request, err := http.NewRequest("POST", "http://example.com?user_id=180&name=tim", strings.NewReader(url.Values{"bodyName": {"isaac"}}.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	var userID int
	var name string
	var bodyName string

	form := NewForm()

	form.AddField("user_id", FieldValue(&IntField{&userID}), QueryParam)
	form.AddField("name", FieldValue(&StringField{&name}), QueryParam)
	form.AddField("bodyName", FieldValue(&StringField{&bodyName}), FormValueParam)

	parseErr := form.Parse(request)

	if parseErr != nil {
		t.Fatal(parseErr)
	}

	if 180 != userID {
		t.Errorf("did not load integer; want %d, got %d", 180, userID)
	}

	if "tim" != name {
		t.Errorf("did not load string; want %q, got %q", "tim", name)
	}

	if "isaac" != bodyName {
		t.Errorf("did not load string; want %q, got %q", "isaac", bodyName)
	}
}
