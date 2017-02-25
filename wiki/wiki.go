package wiki

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Page struct {
	ID       bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name     string
	Title    string
	Author   string
	Contents string
	Tag      []string
	Publish  bool
}

type Wiki struct {
	Page
	Session    *mgo.Session
	DB         string
	Collection string
}

func New(addr string) (*Wiki, error) {
	session, err := mgo.Dial(addr)
	if err != nil {
		return nil, err
	}
	return &Wiki{Session: session,
		DB:         "wiki",
		Collection: "page"}, nil
}

func (w Wiki) CreatePage(page *Page) error {
	c := w.Session.DB(w.DB).C(w.Collection)
	err := c.Insert(page)
	if err != nil {
		return err
	}
	return nil
}

func (w Wiki) ReadPage(name string) (*Page, error) {
	c := w.Session.DB(w.DB).C(w.Collection)
	q := c.Find(bson.M{"name": name})
	n, err := q.Count()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("Page not found")
	}

	page := Page{}
	err = q.One(&page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}
