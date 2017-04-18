package main

import (
	"encoding/json"
	"fmt"
	"github.com/yundream/gowiki/wiki"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
func main() {
	data, err := ioutil.ReadFile("editor.json")
	CheckErr(err)
	page := wiki.Page{}
	err = json.Unmarshal(data, &page)
	CheckErr(err)
	session, err := mgo.Dial("localhost")
	CheckErr(err)
	c := session.DB("wiki").C("page")
	c.Insert(page)
}
