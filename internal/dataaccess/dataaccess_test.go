package dataaccess

import (
	"testing"
	"time"
)

/**
 * Saves Event returns nil | error
 * Test insert
 * Test update
 */
func TestUpSert(t *testing.T) {
	e := Event{}

	// Save Event
	e.Artist = "Greta van Fleet"
	e.Venue = "Forum"
	e.When = time.Now()
	// Test UpSert/insert
	id, _ := e.UpSert()

	if id == "" {
		t.Error("Insert failed - no id")
	}
}

func TestMapFields(t *testing.T) {
	e := Event{Artist: "Stoneroses"}
	m := e.mapFields()

	if e.Artist != m["Artist"] {
		t.Error("Fields are not the same")
	}
}
