package forms

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_New(t *testing.T) {
	var form = New(nil)

	formType := reflect.TypeOf(form)
	if formType != reflect.TypeOf(&Form{}) {
		t.Errorf("expected %v, got %v", reflect.TypeOf(&Form{}), formType)
	}

}

func Test_Has(t *testing.T) {
	var form = New(url.Values{})

	form.Values.Add("foo", "bar")

	val := form.Has("foo")
	if !val {
		t.Errorf("expected %v, got %v", true, val)
	}

	form.Values.Add("baz", "")

	valBaz := form.Has("baz")
	if valBaz {
		t.Errorf("expected %v, got %v", false, val)
	}

}

func Test_Required(t *testing.T) {

	var form = New(url.Values{})

	form.Required("foo")

	if form.Errors.Get("foo") != "This field cannot be blank" {
		t.Errorf("expected %s, got %s", "This field cannot be blank", form.Errors.Get("foo"))
	}
}

func Test_MinLength(t *testing.T) {

	var form = New(url.Values{})

	form.Values.Add("foo", "bar")

	form.MinLength("foo", 5)

	if form.Errors.Get("foo") != "This field must be at least 5 characters long" {
		t.Errorf("expected %s, got %s", "This field must be at least 5 characters long", form.Errors.Get("foo"))
	}

}

func Test_IsEmail(t *testing.T) {

	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{
			name:     "Valid email",
			field:    "foo@bar.com",
			expected: "Invalid email address",
		},
		{
			name:     "Invalid email",
			field:    "bar",
			expected: "Invalid email address",
		},
		{
			name:     "Empty email",
			field:    "",
			expected: "This field cannot be blank",
		},
	}

	for _, tt := range tests {

		var form = New(url.Values{})

		t.Run(tt.name, func(t *testing.T) {

			form.Values.Add("email", tt.field)

			form.IsEmail("email")

			if len(form.Values.Get("email")) == 0 {

				if tt.expected != form.Errors.Get("email") {
					t.Errorf("expected %s, got %s", tt.expected, form.Errors.Get("email"))
				}

			} else {

				if len(form.Errors) > 0 {
					if tt.expected != form.Errors.Get("email") {
						t.Log(form.Errors.Get("email"))
						t.Errorf("expected %s, got %s", tt.expected, form.Errors.Get("email"))
					}

					if form.Values.Get("email") != tt.field {
						t.Errorf("expected %s, got %s", tt.field, form.Values.Get("email"))
					}
				}

			}

			if form.Values.Get("email") != tt.field {
				t.Errorf("expected %s, got %s", tt.field, form.Values.Get("email"))
			}

		})
	}
}

func Test_PasswordsMatch(t *testing.T) {

	tests := []struct {
		name     string
		field1   string
		field2   string
		expected string
	}{
		{
			name:     "Passwords do not match",
			field1:   "foo",
			field2:   "bar",
			expected: "Passwords do not match",
		},
		{
			name:     "Passwords match",
			field1:   "foo",
			field2:   "foo",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var form = New(url.Values{})

			form.Values.Add("password1", tt.field1)
			form.Values.Add("password2", tt.field2)

			form.PasswordsMatch("password1", "password2")

			if tt.expected != form.Errors.Get("password1") {
				t.Errorf("expected %s, got %s", tt.expected, form.Errors.Get("password1"))
			}

		})
	}

}

func Test_Valid(t *testing.T) {

	var form = New(url.Values{})

	form.Values.Add("foo", "bar")

	if !form.Valid() {
		t.Errorf("expected %v, got %v", true, form.Valid())
	}

	form.Required("baz")

	if form.Valid() {
		t.Errorf("expected %v, got %v", false, form.Valid())
	}
}
