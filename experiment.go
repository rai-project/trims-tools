package micro

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Experiment struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at"  bson:"created_at"`
	Hostname  string

	Metadata map[string]string
}
