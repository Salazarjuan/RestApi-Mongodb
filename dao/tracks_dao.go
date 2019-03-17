package dao

import (
	"log"

	. "../models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TracksDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "tracks"
)

func (m *TracksDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
	log.Println("connected to data base! welcome mongo")
}

func (m *TracksDAO) FindAll() ([]Track, error) {
	var tracks []Track
	err := db.C(COLLECTION).Find(bson.M{}).All(&tracks)
	return tracks, err
}

func (m *TracksDAO) FindById(id string) (Track, error) {
	var track Track
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(track)
	return track, err
}

func (m *TracksDAO) Insert(track Track) error {
	err := db.C(COLLECTION).Insert(&track)
	return err
}

func (m *TracksDAO) Delete(track Track) error {
	err := db.C(COLLECTION).Remove(&track)
	return err
}

func (m *TracksDAO) Update(track Track) error {
	err := db.C(COLLECTION).UpdateId(track.ID, &track)
	return err
}
