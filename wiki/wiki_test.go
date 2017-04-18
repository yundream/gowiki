package wiki

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var wiki *Wiki

func init() {
	var err error
	wiki, err = New("localhost")
	if err != nil {
		log.Fatal(err)
	}
}

func Test_CreatePage(t *testing.T) {
	page := Page{
		Name:     "TestPage",
		Title:    "FrontPage",
		Contents: "Hello World",
		Tag:      []string{"joinc", "news"},
		Publish:  true}
	err := wiki.CreatePage(&page)
	assert.Nil(t, err, "")
}

func Test_ReadPage(t *testing.T) {
	page, err := wiki.ReadPage("TestPage")
	assert.Nil(t, err, "")
	assert.Equal(t, page.Name, "TestPage")
	assert.Equal(t, page.Title, "FrontPage")
}

func Test_EditPage(t *testing.T) {
	page, err := wiki.ReadPage("editor")
	assert.Nil(t, err, "")
	t.Log(page.Contents)
}
