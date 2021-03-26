package handlers 

import (
	"testing"
	"time"
)
//To be valid name of the function should start with Test
func TestHumanDate(t *testing.T) {
	//initializing and passing a new time.Time object and pas it to
	//humanDate function
	tm := time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC)
	hd := HumanDate(tm)

	if hd != "17 Dec 2020 at 10:00" {
		t.Errorf("want %q; got %q", "17 Dec 2020 at 10:00", hd)
	}
}
