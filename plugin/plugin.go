package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	fmt.Println("Plugins ", plugins)
	return plugins, nil
}

type opt struct {
	Name string
	Age  int
}

func (p PlugIns) Exec(fname string, fdata string) (string, error) {
	if sym, ok := p.pluginList["Function_"+fname]; ok {
		r := sym.(func(string, string) string)(fdata, "yundream")
		return r, nil
	} else {
		fmt.Println("Plugin ", fname, "Exec Error")
	}
	return "", errors.New("ERROR")
}

func (p PlugIns) List() string {
	return ""
}
