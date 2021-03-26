package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

//Custom Form struct, which anonymously embeds a url.Values object
//to hold the form data and Errors field to hold any validation errors
//for the form data.

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9.!]")

type Form struct {
	url.Values
	Errors errors
}

//New() function to initialize a custom Form struct
func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

//MinLength checks if the specific field contains a minimum number of
//characters.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is to short (minumum is %d characters long)", d))
	}
}

//MatchesPattern method to check tha a specific filed in the form matches
//a regular expression. If the chech fails then add the
//appropriate messsage the the form errors
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

//Required field to check that specific fields in the form data present
//and not blank. If any field fail this check, add the appropriate
//message to the form error
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field) //Get is the method of the anonymous field url
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

//MaxLength method to check if that a specific field in the form
//contains a maxium number of characters.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("this is field is too long maximum is %d characters long", d))
	}
}

//PermittedValues checks if that a specific field in the form matches one
//of a set of specific permitted values

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}

	}
	f.Errors.Add(field, "This field is invalid")
}

//Valid method returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//regexp.MustCompile() to parse a pattern and compile a regular expression
//for sanity checking the format of an email address
//this returns *regexp.Regexp object, or panics in the event of an error
//we do this once at runtime, store the compiled regular expression
//object in a variable
