package wiki

import (
	"gopkg.in/mgo.v2"
)

type Page struct {
	Title    string
	Author   string
	Contents string
	Tag      []string
	Publish  bool
}

type Wiki struct {
	Page
	D *mgo.Session
}

func New(addr string) (*Wiki, error) {
	db, err := mgo.Dial(addr)
	if err != nil {
		return nil, err
	}
	return Wiki{D: db}, nil
}

func (w *Wiki) Create() error {
	return nil
}

func (w *Wiki) ReadFromID(id string) error {
	return nil
}

func (w *Wiki) ReadFromName(name string) error {
	return nil
}
