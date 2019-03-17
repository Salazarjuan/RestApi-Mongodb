package models

import "gopkg.in/mgo.v2/bson"

type Track struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Author     string        `bson:"author" json:"author"`
	Album      string        `bson:"album" json:"album"`
	CoverImage string        `bson:"cover_image" json:"cover_image"`
}
