package dataaccess

import (
	"encoding/json"
	"errors"
	"time"

	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

var db *c.DB

/**
 *
 */
type Event struct {
	Artist string
	Venue  string
	When   time.Time
	Id     string
}

func init() {
	db, _ = c.Open("./clover.db")
	// Create collection
	if _, err := db.HasCollection("events"); err != nil {
		db.CreateCollection("events")
	}
}

func (e *Event) UpSert() (string, error) {
	// If Id is "" insert
	// Otherwise update
	doc := d.NewDocument()
	var err error
	if e.Id == "" {
		doc.Set("artist", e.Artist)
		doc.Set("venue", e.Venue)
		e.Id, err = db.InsertOne("events", doc)
		return e.Id, err
	} else {
		db.Update(query.NewQuery("events").Where(query.Field("Id").Eq(e.Id)), e.mapFields())
		return e.Id, err
	}
	db.Close()

	return "", errors.New("UpSert failed")
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

/**
 * Save batch of Events
 */
func (e *Event) SaveMultiple([]Event) error {
	return errors.New("Failed")
}

/**
 * Get Single Event
 **/
func (e *Event) GetEvent(id int64) (Event, error) {
	return Event{}, errors.New("Failed")
}

/**
 * Returns all Events
 */
func (e *Event) GetAllEvents() ([]Event, error) {
	return []Event{}, errors.New("Failed")
}
