package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

type Form struct {
	Values url.Values
	Errors errors
}

func New(form url.Values) *Form {
	return &Form{
		Values: form,
		Errors: errors{},
	}
}

func (f *Form) Has(field string) bool {
	x := f.Values.Get(field)
	if x == "" {
		return false
	}
	return true
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Values.Get(field)
		if value == "" {
			f.Errors[field] = append(f.Errors[field], "This field cannot be blank")
		}
	}
}

func (f *Form) MinLength(field string, length int) {
	x := f.Values.Get(field)
	if len(x) < length {
		f.Errors[field] = append(f.Errors[field], fmt.Sprintf("This field must be at least %d characters long", length))
	}
}

func (f *Form) IsEmail(field string) {
	if len(f.Values.Get(field)) == 0 {
		f.Errors[field] = append(f.Errors[field], "This field cannot be blank")
		return
	}
	emailAdrress, err := mail.ParseAddress(f.Values.Get(field))

	if err != nil {
		f.Errors[field] = append(f.Errors[field], "Invalid email address")
	} else {
		f.Values.Set(field, emailAdrress.Address)
	}

}

func (f *Form) PasswordsMatch(field1 string, field2 string) {
	if !strings.EqualFold(f.Values.Get(field1), f.Values.Get(field2)) {
		f.Errors[field1] = append(f.Errors[field1], "Passwords do not match")
		f.Errors[field2] = append(f.Errors[field2], "Passwords do not match")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
