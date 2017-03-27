package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"plugin"
)

type Fmap map[string]plugin.Symbol

type PlugIns struct {
	pluginList Fmap
}

func Load() (*PlugIns, error) {
	plugins := &PlugIns{}
	plugins.pluginList = make(Fmap)

	files, err := ioutil.ReadDir("./plugin")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		sofile := fmt.Sprintf("./plugin/%s/%s.so", file.Name(), file.Name())
		fmt.Println("Plugin loading...", sofile)
		if _, err := os.Stat(sofile); os.IsNotExist(err) {
			fmt.Println(err.Error())
			continue
		}
		p, err := plugin.Open(sofile)
		if err != nil {
			fmt.Println("Loading Error ", err.Error())
			return nil, err
		}
		sym, err := p.Lookup("Function_" + file.Name())
		if err != nil {
			return nil, err
		}
		plugins.pluginList["Function_"+file.Name()] = sym
	}
	return plugins, nil
}

type opt struct {
	Name string
	Age  int
}

func (p PlugIns) Exec(fname string, fdata string, w http.ResponseWriter, r *http.Request) (string, error) {
	if sym, ok := p.pluginList["Function_"+fname]; ok {
		r := sym.(func(string, string, http.ResponseWriter, *http.Request) string)(fdata, "yundream", w, r)
		return r, nil
	} else {
	}
	return "", errors.New("ERROR")
}

func (p PlugIns) List() string {
	return ""
}
