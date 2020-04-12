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

// RouteHandler returns a route handler struct
type RouteHandler struct {
	Db Database
}

// Database is a struct containing a bolt db
type Database struct {
	Db *bolt.DB
}

// User is a struct that represents a user
type User struct {
	ID                int
	Username          string
	Name              string
	Zipcode           string
	Groups            []Group
	RestaurantsTried  []string
	RestaurantsMissed []string
}

// Group is a struct that represents a group
type Group struct {
	ID                int      `json:"groupID"`
	GroupName         string   `json:"groupName,omitempty"`
	Members           []int   `json:"groupMembers,omitempty"`
	RestaurantsTried  []int `json:"restaurantsTried,omitempty"`
	RestarantsMissed  []int `json:"restaurantsMissed,omitempty"`
	CollectiveZipcode string   `json:"collectiveZip,omitempty"`
}

// Restaurant is a struct that represents a restaurant
type Restaurant struct {
	ID            int
	YelpId        string
	Liked         bool
	MichelinStars float64
}

type bucketType int

const (
	BktGroup bucketType = iota
	BktUser
	BktRestaurant
)

func (bt bucketType) String() string {
	return [...]string{"group", "user", "restaurant"}[bt]
}

// NewRouter opens a database, and retruns a Database struct
func NewRouter() (*RouteHandler, error) {
	db, err := bolt.Open("/Users/Spyro/Developer/go/src/db/wfl.db", 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("NewRouter: bolt.Open: %v", err)
	}
	defer db.Close()
	database := &Database{Db: db}
	return &RouteHandler{Db: *database}, nil
}

// CreateGroup adds a group to the database groups bucket
func (db *Database) CreateGroup(g *Group) (int, error) {
	var ID int
	err := db.Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BktGroup.String()))
		if err != nil {
			return fmt.Errorf("CreateGroup: CreateBucket: %v", err)
		}

		id, _ := b.NextSequence()
		g.ID = int(id)
		ID = int(id)

		buf, err := json.Marshal(g)
		if err != nil {
			return fmt.Errorf("CreateGroup: json.marshal: %v", err)
		}
		return b.Put(itob(g.ID), buf)
	})
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: db.Db.Update: %v", err)
	}
	return ID, nil
}

// CreateUser adds a user to the user bucket in the database
func (db *Database) CreateUser(u *User) (int, error) {
	var ID int
	err := db.Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BktUser.String()))
		if err != nil {
			return fmt.Errorf("CreateUser: CreateBucket: %v", err)
		}

		id, _ := b.NextSequence()
		u.ID = int(id)
		ID = int(id)

		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return b.Put(itob(u.ID), buf)
	})
	if err != nil {
		return 0, err
	}
	return ID, nil
}

// CreateRestaurant adds a restarant to the restarant bucket in the database
func (db *Database) CreateRestaurant(r *Restaurant) (int, error) {
	var ID int
	err := db.Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BktRestaurant.String()))
		if err != nil {
			return fmt.Errorf("CreateRestaurant: CreateBucket: %v", err)
		}

		id, _ := b.NextSequence()
		r.ID = int(id)
		ID = int(id)

		buf, err := json.Marshal(r)
		if err != nil {
			return err
		}
		return b.Put(itob(r.ID), buf)
	})
	if err != nil {
		return 0, nil
	}
	return ID, nil
}

// QueryFromDb querys the given group in the database for the entry with a specific id
func (db *Database) QueryFromDb(id int, bucket bucketType) ([]byte, error) {
	var entry []byte
	err := db.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.String()))
		if b == nil {
			return fmt.Errorf("Bucket %v doesn't exist", bucket.String())
		}
		entry = b.Get(itob(id))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("QueryFromDb: db.Db.View: %v", err)
	}
	return entry, nil
}

// DeleteFromDb deletes an entry from the database in the given bucket
func (db *Database) DeleteFromDb(id int, bucket bucketType) error {
	return db.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.String()))
		return b.Delete(itob(id))
	})
}

// HandlerCreateGroup is used by the RouteHandler to add a group to the database.
// This returns the ID of the created group
func (rh *RouteHandler) HandlerCreateGroup(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var g Group
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("HandlerCreateGroup: json.decode: %v", err)
	}

	ID, err := rh.Db.CreateGroup(&g)
	if err != nil {
		return nil, fmt.Errorf("HandlerCreateGroup: Db.CreateGroup: %v", err)
	}

	b, err := json.Marshal(Group{ID: ID})
	if err != nil {
		return nil, fmt.Errorf("HandlerCreateGroup: json.marshal: %v", err)
	}
	return b, nil
}

// HandlerQueryDb is used by the RouteHandler to Query the a group from the url scheme
func (rh *RouteHandler) HandlerQueryDb(w http.ResponseWriter, r *http.Request, idStr string, bucket bucketType) ([]byte, error) {
	varID := mux.Vars(r)[idStr]
	id, err := strconv.Atoi(varID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("HandlerQueryDb: strconv.Atio: %v", err)
	}
	b, err := rh.Db.QueryFromDb(id, bucket)
	if err != nil {
		return nil, fmt.Errorf("HandlerQueryDb: QueryFromDb: %v", err)
	}
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, nil
	}
	return b, err
}

// HandlerDeleteDb is used by RouteHandler to delete an entry from the bucket in the db.
func (rh *RouteHandler) HandlerDeleteDb(w http.ResponseWriter, r *http.Request, idStr string, bucket bucketType) error {
	varID := mux.Vars(r)[idStr]
	id, err := strconv.Atoi(varID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("HandlerDeleteDb: error %v", err)
	}
	return rh.Db.DeleteFromDb(id, bucket)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
