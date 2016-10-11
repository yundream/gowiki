# gowiki
## Plugin 시스템
```go
package plugin

import ()

type Plugin struct {
    version string
}

func (p Plugin) Version() string {
    return p.version
}
func (p *Plugin) SetVersion(ver string) {
    p.version = ver
}
```

```go
# plugin/helloworld
package helloworld

import (
    "bitbucket.org/dream_yun/testimport/plugin"
)

// plugin을 embeded 한다.
type Hello struct {
    plugin.Plugin
}
func New(p plugin.Plugin) *Hello {
    return &Hello{p}
}

// 플러그인 함수를 개발한다.
func (h Hello) Sum(a, b int) int {
    return a + b
}
```

```go
package main

import (
    "bitbucket.org/dream_yun/testimport/plugin"
    "bitbucket.org/dream_yun/testimport/plugin/helloworld"
    "fmt"
    "reflect"
)

func main() {
    p := plugin.Plugin{}
    p.SetVersion("v3.0")

    // plugin을 읽어온다.
    // 실제 코드에서는 plugin/ 디렉토리에서 읽어서 처리하게 한다.
    plugin_helloworld := helloworld.New(p)
    method := reflect.TypeOf(plugin_helloworld)

    F := make(map[string]reflect.Value)
    for i := 0; i < method.NumMethod(); i++ {
        name := method.Method(i).Name
        fmt.Println("Reflect method ", name)
        reflectMethd := reflect.ValueOf(plugin_helloworld).MethodByName(name)
        F["helloworld/"+method.Method(i).Name] = reflectMethd

    }
    mm, ok := F["helloworld/Sum"]
    if ok {
        r := mm.Call([]reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)})
        fmt.Println(r[0].Int())
    }
}
```
