package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
)

type RouteHandler struct {
	Db Database
}

// Database is a struct containing a bolt db
type Database struct {
	Db *bolt.DB
}

// User is a struct that represents a user
type User struct {
	Username          string
	Name              string
	Zipcode           string
	Groups            []Group
	RestaurantsTried  []string
	RestaurantsMissed []string
}

// Group is a struct that represents a group
type Group struct {
	ID         int   `json:"groupID"`
	GroupName         string   `json:"groupName"`
	Gembers           []User   `json:"groupMembers"`
	RestaurantsTried  []string `json:"restaurantsTried,omitempty"`
	RestarantsMissed  []string `json:"restaurantsMissed,omitempty"`
	CollectiveZipcode string   `json:"collectiveZip,omitempty"`
}

// Restaurant is a struct that represents a restaurant
type Restaurant struct {
	YelpID        string
	Liked         bool
	MichelinStars float64
}

// NewStore opens a database, and retruns a Database struct
func NewRouter() (*RouteHandler, error) {
	db, err := bolt.Open("/Users/Spyro/Developer/go/src/db/wfl.db", 0666, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	database := &Database{Db: db}
	return &RouteHandler{Db: *database}, nil
}

func (db *Database) CreateGroup(g *Group) error {
	return db.Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("groups"))
		if err != nil {
			return fmt.Errorf("Write: CreateBucket: %v", err)
		}

		id, _ := b.NextSequence()
		g.ID = int(id)

		buf, err := json.Marshal(g)
		if err != nil {
			return err
		}
		
		return b.Put(itob(g.ID), buf)
	})
}

func (db *Database) QueryGroup(id int) ([]byte) {
	var group []byte
	db.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("groups"))
		group = b.Get(itob(id))
		return nil
	})
	return group
}

func (h *RouteHandler) GroupQueryHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	groupID := mux.Vars(r)["groupId"]
	id, err := strconv.Atoi(groupID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	b := h.Db.QueryGroup(id)
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, nil
	}
	return b, err
}

func (h *RouteHandler) GroupCreateHandler(w http.ResponseWriter, r *http.Request) (error) {
	var g Group
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		return err
	}
	return h.Db.CreateGroup(&g)
}

func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}