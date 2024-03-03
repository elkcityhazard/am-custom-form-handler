package forms

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_Add(t *testing.T) {
	var form = New(nil)

	err := fmt.Errorf("bar")

	form.Errors.Add("foo", err.Error())

	if form.Errors.Get("foo") != err.Error() {
		t.Errorf("expected %s, got %s", err, form.Errors.Get("foo"))
	}
}

func Test_Get(t *testing.T) {

	var vals = url.Values{}
	var errors = errors{}

	var form = New(vals)
	errors.Add("baz", "")

	result := form.Errors.Get("baz")
	if result != "" {
		t.Errorf("expected %s, got %s", "", result)
	}

	errors.Add("foo", "bar")

	form.Errors = errors

	if form.Errors.Get("foo") != "bar" {
		t.Errorf("expected %s, got %s", "bar", form.Errors.Get("foo"))
	}

}
