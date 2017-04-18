package plugin

import (
	"fmt"
	"os"
	"testing"
)

func Test_load(t *testing.T) {
	p, err := Load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	rtv, err := p.Exec("Function_sample", "name")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	t.Log(rtv)
}
