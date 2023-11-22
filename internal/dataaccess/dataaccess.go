// Package dataaccess implements methods for creating,
// updating and fetching data objects such as events
package dataaccess

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

const (
	collectionName = "events"
)

var db *c.DB

// An Event holds relevant fields
type Event struct {
	Artist string
	Venue  string
	When   string
	ID     string
}

func init() {
	var err error
	file, _ := filepath.Abs(os.Getenv("GLOVER_PATH"))
	db, err = c.Open(file)

	// Handle open err
	if err != nil {
		log.Fatal("Clover could not open ", err)
	}
	coll, _ := db.HasCollection(collectionName)
	// Create collection
	if coll == false {
		log.Println("Create collection")
		db.CreateCollection(collectionName)
	}
}

// UpSert inserts or updates the current Event
// and returns either a string ID or an error
func (e *Event) UpSert() (string, error) {
	// If Id is "" insert
	// Otherwise update
	doc := d.NewDocument()
	var err error
	if e.ID == "" {
		doc.Set("Artist", e.Artist)
		doc.Set("Venue", e.Venue)
		e.ID, err = db.InsertOne(collectionName, doc)
		return e.ID, err
	}
	err = db.Update(query.NewQuery(collectionName).Where(query.Field("ID").Eq(e.ID)), e.mapFields())
	return e.ID, err
}

/**
 * TODO: Add fields to map
 */
func (e *Event) mapFields() map[string]interface{} {
	var event = make(map[string]interface{})
	d, _ := json.Marshal(e)
	json.Unmarshal(d, &event)

	return event
}

// SaveMultiple saves the events passed
// and returns an error or nil
func (e *Event) SaveMultiple([]Event) error {
	return errors.New("Failed")
}

// GetEvent returns the Event found based on the passed ID
// or an error
func (e *Event) GetEvent(id string) (Event, error) {
	doc, err := db.FindById(collectionName, id)

	if err != nil {
		return Event{}, err
	}

	return Event{Artist: doc.Get("Artist").(string), Venue: doc.Get("Venue").(string)}, nil
}

// GetAllEvents returns all the events in the collection
func (e *Event) GetAllEvents() ([]Event, error) {
	var events []Event
	docs, err := db.FindAll(query.NewQuery(collectionName).Sort())

	for _, doc := range docs {
		events = append(events, Event{
			Artist: doc.Get("Artist").(string),
			Venue:  doc.Get("Venue").(string),
		})
	}
	return events, err
}
