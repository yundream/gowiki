package wiki

import (
	"bufio"
	"bytes"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
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
	compiler   *Compiler
}

func New(addr string) (*Wiki, error) {
	session, err := mgo.Dial(addr)
	if err != nil {
		return nil, err
	}
	compiler := Compiler{}
	err = compiler.LoadPlugin()
	if err != nil {
		return nil, err
	}
	return &Wiki{Session: session,
		DB:         "wiki",
		Collection: "page",
		compiler:   &compiler,
	}, nil
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
	var buffer bytes.Buffer
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

	compilerIns := w.compiler.NewIns()
	compilerIns.Start("myid")

	scanner := bufio.NewReader(strings.NewReader(page.Contents))
	for {
		line, _, err := scanner.ReadLine()
		if err != nil {
			break
		}
		compilerIns.Line(string(line))
		switch compilerIns.TextType {
		case PROCESSOR_OPEN:
			if compilerIns.Processor() == PROCESSOR_CLOSE {
				buffer.WriteString(compilerIns.String())
			}
		default:
			compilerIns.List().Head().EscapeString().Body()
			buffer.WriteString(compilerIns.String())
		}
	}
	page.Contents = buffer.String()
	return &page, nil
}
