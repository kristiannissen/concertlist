package dataaccess

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setup() {
	// Delete collection before tests are run
	// db.DropCollection(collectionName)

}

func teardown() {
	path, _ := os.Getwd()
	file, _ := filepath.Abs(strings.Join([]string{path, "events.json"}, string(separator)))
	db.ExportCollection(collectionName, file)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

/**
 * Saves Event returns nil | error
 * Test insert
 * Test update
 */
func TestUpSert(t *testing.T) {
	e := Event{Artist: "Rival Sons"}
	id, _ := e.UpSert()

	t.Logf("id %s", id)

	if id == "" {
		t.Error("Event Id not set")
	}

}

func TestMapFields(t *testing.T) {
	e := Event{Artist: "Stoneroses"}
	m := e.mapFields()

	if e.Artist != m["Artist"] {
		t.Error("Fields are not the same")
	}
}

func TestGetAllEvents(t *testing.T) {
	e := Event{Artist: "Led Zeppelin"}
	id, err := e.UpSert()
	t.Log("id ", id, err)

	events, err := e.GetAllEvents()
	t.Log("events ", events)

	if err != nil {
		t.Error("Could not get all events")
	}
}
