package greetings


import (
	"testing"
	"regexp"
)


func TestHelloName(t *testing.T) {
	name := "Daniel"
	want := regexp.MustCompile(`\b`+name+`\b`)
	msg, err := Hello("Daniel")
	if !want.MatchString(msg) || err != nil {
		t.Errorf(`Hello("Daniel") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}


func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Errorf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
