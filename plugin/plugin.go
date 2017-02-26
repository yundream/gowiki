package plugin

import (
	"errors"
	"fmt"
	"github.com/yundream/gowiki/handler"
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

	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		func(os.FileInfo) {
			if !file.IsDir() {
				return
			}
			so := fmt.Sprintf("./%s/%s.so", file.Name(), file.Name())
			if _, err := os.Stat(so); os.IsNotExist(err) {
				return
			}
			fmt.Println("Plugin loading...", so)
		}(file)
	}

	p, err := plugin.Open("./sample/sample.so")
	if err != nil {
		fmt.Println("Loading Error ", err.Error())
		return nil, err
	}
	sym, err := p.Lookup("Function_sample")
	if err != nil {
		return nil, err
	}
	plugins.pluginList["Function_sample"] = sym

	return plugins, nil
}

func (p PlugIns) Exec(name string, data string) (string, error) {
	if sym, ok := p.pluginList[name]; ok {
		r := sym.(func(string, handler.Options) string)(data, handler.Options{Name: "yundream"})
		return r, nil
	}
	return "", errors.New("ERROR")
}

func (p PlugIns) List() string {
	return ""
}
