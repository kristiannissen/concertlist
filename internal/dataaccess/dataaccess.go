package dataaccess

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

const (
	collectionName = "events"
	separator      = os.PathSeparator
	storage        = "clover.db"
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
	var err error
	path, _ := os.Getwd()
	abs, _ := filepath.Abs(strings.Join([]string{path, storage}, string(separator)))
	db, err = c.Open(abs)
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

func (e *Event) UpSert() (string, error) {
	// If Id is "" insert
	// Otherwise update
	doc := d.NewDocument()
	var err error
	if e.Id == "" {
		doc.Set("Artist", e.Artist)
		doc.Set("Venue", e.Venue)
		e.Id, err = db.InsertOne(collectionName, doc)
		return e.Id, err
	} else {
		err = db.Update(query.NewQuery(collectionName).Where(query.Field("Id").Eq(e.Id)), e.mapFields())
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
func (e *Event) GetEvent(id string) (Event, error) {
	return Event{}, errors.New("Failed")
}

/**
 * Returns all Events
 */
func (e *Event) GetAllEvents() ([]Event, error) {
	var events []Event
	docs, err := db.FindAll(query.NewQuery(collectionName).Sort())
	// defer db.Close()

	log.Println(docs)

	for _, doc := range docs {
		events = append(events, Event{Artist: doc.Get("Artist").(string)})
	}
	return events, err
}
